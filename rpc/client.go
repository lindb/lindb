package rpc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/rpc/pkg/batch"
)

// the rpc implementation heavily refer to tidb https://github.com/pingcap/tidb

// rpc config
const (
	dialTimeout               = 5 * time.Second
	keepAlive                 = 10 * time.Second
	keepAliveTimeout          = 3 * time.Second
	maxConnSize               = 1
	maxBatchSize              = 10
	grpcInitialWindowSize     = 1 << 30
	grpcInitialConnWindowSize = 1 << 30
)

// BaseClient provides a unify stream call interface, different requests can be converted to batchRequests
// and transfer to server in a stream way.
type BaseClient interface {
	// Close releases connection resources
	Close()
	// GetConn returns a grpc.ClientConn randomly which can be used to construct grpcServiceClient and send single request
	GetConn(address string) (*grpc.ClientConn, error)
	// SendStreamRequest batch requests to send in a stream way
	SendStreamRequest(cxt context.Context, address string, request *Request, timeout time.Duration) (*Response, error)
}

// baseClient implements BaseClient handling connection pool, stream requests, timeout checks
type baseClient struct {
	sync.RWMutex
	isClosed bool
	// address(host+port) -> conn pool(grpc connections)
	conns map[string]*connPool
}

// NewBaseClient returns a new BaseClient implementation
func NewBaseClient() BaseClient {
	return &baseClient{
		conns: make(map[string]*connPool),
	}
}

// getCoonPool returns a connPool for given address
func (c *baseClient) getCoonPool(address string) (*connPool, error) {
	c.RLock()
	if c.isClosed {
		c.RUnlock()
		return nil, errors.Errorf("baseClient is closed")
	}
	pool, ok := c.conns[address]
	c.RUnlock()
	if !ok {
		var err error
		pool, err = c.createCoonPool(address)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}

// createCoonPool creates a connPool for given address, goroutine safe
func (c *baseClient) createCoonPool(address string) (*connPool, error) {
	c.Lock()
	defer c.Unlock()
	pool, ok := c.conns[address]
	if !ok {
		var err error
		pool, err = newConnPool(address, maxConnSize, maxBatchSize)
		if err != nil {
			return nil, err
		}
		c.conns[address] = pool
	}
	return pool, nil
}

// GetConn returns a underlying grpc.ClientConn randomly
func (c *baseClient) GetConn(address string) (*grpc.ClientConn, error) {
	pool, err := c.getCoonPool(address)
	if err != nil {
		return nil, err
	}
	return pool.GetCoon(), nil
}

// SendStreamRequest sends request and wait for response, the request will be batched and send in a stream way
func (c *baseClient) SendStreamRequest(
	ctx context.Context,
	address string,
	req *Request,
	timeout time.Duration) (*Response, error) {
	pool, err := c.getCoonPool(address)
	if err != nil {
		return nil, err
	}

	batchReq := req.ToBatchRequest()
	if batchReq == nil {
		return nil, errors.Errorf("unknown request type")
	}

	// entry wrappers the request, and constructs a response channel, the send goroutine wait on the channel for response
	entry := &batchRequestEntry{
		req: batchReq,
		res: make(chan *batch.BatchResponse_Response, 1),
	}

	ctx1, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case pool.batchRequestsCh <- entry:
	case <-ctx1.Done():
		return nil, errors.Errorf("send request timeout to %s", address)
	}

	logger.GetLogger().Debug("entry send to channel, wait for resp")

	select {
	case res, ok := <-entry.res:
		// channel closed
		if !ok {
			return nil, entry.err
		}
		return FromBatchResponse(res), nil
	case <-ctx1.Done():
		atomic.StoreInt32(&entry.canceled, 1)
		return nil, errors.Errorf("send request canceled or timeout")
	}
}

// closeConns closes all conn pools
func (c *baseClient) closeConns() {
	c.Lock()
	defer c.Unlock()
	if !c.isClosed {
		c.isClosed = true

		for _, pool := range c.conns {
			pool.Close()
		}
	}
}

// Close closes the baseClient, releases connections
func (c *baseClient) Close() {
	c.closeConns()
}

// parallel connections pool to the same address for efficiency
type connPool struct {
	address string

	conns []*grpc.ClientConn

	// a goroutine check request timeout in background
	leasesCh chan *lease

	// batchRequestCh buffers the requests and sends in a batch
	batchRequestsCh chan *batchRequestEntry
	// stream baseClient based on grpc.ClientConn
	batchRequestsClients []*batchRequestsClient

	// atomic index to round robbin requests to conns
	roundRobbinIndex uint32
}

// newConnPool returns a connPool for the target address, with maxConnSize, maxBatchSize constrains
func newConnPool(address string, maxConnSize uint32, maxBatchSize uint32) (*connPool, error) {
	pool := &connPool{
		address:              address,
		conns:                make([]*grpc.ClientConn, maxConnSize),
		leasesCh:             make(chan *lease, 1024),
		batchRequestsCh:      make(chan *batchRequestEntry, maxBatchSize),
		batchRequestsClients: make([]*batchRequestsClient, 0, maxConnSize),
	}

	if err := pool.init(); err != nil {
		return nil, err
	}
	return pool, nil
}

// init establishes the connection and start up batchSend, batchRecv, checkTimeout go routines
func (pool *connPool) init() error {
	for i := range pool.conns {
		ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
		// connect
		conn, err := grpc.DialContext(
			ctx,
			pool.address,
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithBackoffMaxDelay(time.Second*3),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                keepAlive,
				Timeout:             keepAliveTimeout,
				PermitWithoutStream: true,
			}),
		)

		cancel()

		if err != nil {
			pool.Close()
			logger.GetLogger().Error("init conn pool error", zap.Stack("stack"), zap.Error(err))
			return err
		}

		pool.conns[i] = conn

		// get stream client
		gCli := batch.NewBatchServiceClient(conn)
		gStreamCli, err := gCli.StreamBatchRequest(context.TODO())
		if err != nil {
			pool.Close()
			logger.GetLogger().Error("create stream baseClient err", zap.Stack("stack"), zap.Error(err))
			return err
		}

		batchCli := &batchRequestsClient{
			address: pool.address,
			conn:    conn,
			client:  gStreamCli,
		}
		pool.batchRequestsClients = append(pool.batchRequestsClients, batchCli)

		go batchCli.batchReceiveLoop()
	}

	go checkLeasesTimeoutLoop(pool.leasesCh)

	go pool.batchSendLoop()

	return nil
}

func (pool *connPool) batchSendLoop() {
	defer func() {
		// recover to avoid the client panic
		if r := recover(); r != nil {
			logger.GetLogger().Error("recover from batchSendLoop", zap.Stack("stack"), zap.Any("recover", r))
			go pool.batchSendLoop()
		}
	}()

	entries := make([]*batchRequestEntry, 0, maxBatchSize)
	requests := make([]*batch.BatchRequest_Request, 0, maxBatchSize)
	requestIDs := make([]uint64, 0, maxBatchSize)

	for {
		// choose a conn by round-robbin
		next := atomic.AddUint32(&pool.roundRobbinIndex, 1) % uint32(len(pool.conns))

		batchRequestCli := pool.batchRequestsClients[next]

		entries = entries[:0]
		requests = requests[:0]
		requestIDs = requestIDs[:0]

		pool.fetchAllPendingRequests(&entries, &requests)
		// todo check server payload, wait a bit more to collect more requests

		// connection closed
		if len(entries) == 0 {
			logger.GetLogger().Warn("pending requests is empty, coon is closed, batchSendLoop return")
			return
		}

		logger.GetLogger().Debug("fetch entries", zap.Int("size", len(entries)))

		length := removeCanceledRequests(&entries, &requests)

		logger.GetLogger().Debug("after remove canceled entries", zap.Int("size", length))

		// all canceled
		if length == 0 {
			continue
		}

		// assign requestID
		maxBatchID := atomic.AddUint64(&batchRequestCli.idSeq, uint64(length))
		batchIDBeg := maxBatchID - uint64(length)
		for i := 0; i < length; i++ {
			requestID := uint64(i) + batchIDBeg
			requestIDs = append(requestIDs, requestID)
		}

		batchRequest := &batch.BatchRequest{
			RequestIDs: requestIDs,
			Requests:   requests,
		}

		// protect the error handling in BatchRecv from changing service baseClient
		batchRequestCli.clientLock.Lock()
		for i, requestID := range batchRequest.RequestIDs {
			batchRequestCli.batched.Store(requestID, entries[i])
		}
		err := batchRequestCli.client.Send(batchRequest)
		batchRequestCli.clientLock.Unlock()

		if err != nil {
			logger.GetLogger().Error("batch requests send error", zap.Stack("stack"), zap.Error(err))
			batchRequestCli.failPendingRequests(err)
		}
	}

}

// collect requests then send in a stream request
func (pool *connPool) fetchAllPendingRequests(
	entries *[]*batchRequestEntry,
	requests *[]*batch.BatchRequest_Request) {

	// Block on the first element
	entry, more := <-pool.batchRequestsCh

	if more {
		*entries = append(*entries, entry)
		*requests = append(*requests, entry.req)
	} else {
		logger.GetLogger().Info("batch request channel is closed")
		return
	}

	// fetch more if possible
	for len(*entries) < maxBatchSize {
		select {
		case entry = <-pool.batchRequestsCh:
			*entries = append(*entries, entry)
			*requests = append(*requests, entry.req)
		default:
			return
		}
	}
}

// remove canceled requests from input entries and requests
func removeCanceledRequests(
	entries *[]*batchRequestEntry,
	requests *[]*batch.BatchRequest_Request) int {
	validateEntries := (*entries)[:0]
	validateRequests := (*requests)[:0]
	for _, entry := range *entries {
		if !entry.isCanceled() {
			validateEntries = append(validateEntries, entry)
			validateRequests = append(validateRequests, entry.req)
		} else {
			logger.GetLogger().Warn("entry canceled", zap.Any("entry", entry.req))
		}
	}

	*entries = validateEntries
	*requests = validateRequests
	return len(*entries)
}

// GetCoon returns ClientConn roundRobbin
func (pool *connPool) GetCoon() *grpc.ClientConn {
	index := atomic.AddUint32(&pool.roundRobbinIndex, 1) % uint32(len(pool.conns))
	return pool.conns[index]
}

// Close closes pooled connections
func (pool *connPool) Close() {
	// close batchReceiveLoop
	for _, c := range pool.batchRequestsClients {
		atomic.StoreInt32(&c.closed, 1)
	}
	close(pool.batchRequestsCh)

	for i, coon := range pool.conns {
		if coon != nil {
			err := coon.Close()
			if err != nil {
				logger.GetLogger().Error("coon pool close error", zap.Stack("stack"), zap.Error(err))
			}
			pool.conns[i] = nil
		}
	}

	close(pool.leasesCh)

}

// lease holds a deadline and a CancelFunc, a goroutine will check the Leases, call CancelFunc after deadline
type lease struct {
	Cancel context.CancelFunc
	// time.UnixNano value
	deadline int64
}

// checkLeasesTimeoutLoop checks lease deadline, call with goroutine
func checkLeasesTimeoutLoop(ch <-chan *lease) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	array := make([]*lease, 0, 1024)

	for {
		select {
		case item, ok := <-ch:
			// ch closed
			if !ok {
				return
			}
			array = append(array, item)
		case now := <-ticker.C:
			array = keepOnlyActive(array, now.UnixNano())

		}
	}
}

// keepOnlyActive removes lease with deadline before or equal than now
func keepOnlyActive(array []*lease, now int64) []*lease {
	index := 0
	for i, item := range array {
		deadline := atomic.LoadInt64(&item.deadline)
		if deadline == 0 || deadline > now {
			array[index] = array[i]
			index++
		} else {
			item.Cancel()
		}
	}
	return array[:index]
}

// batchRequestEntry includes a batchRequest and a channel for corresponding response
type batchRequestEntry struct {
	req *batch.BatchRequest_Request
	res chan *batch.BatchResponse_Response

	canceled int32
	err      error
}

func (b *batchRequestEntry) isCanceled() bool {
	return atomic.LoadInt32(&b.canceled) == 1
}

// batchRequestsClient wrappers grpc stream client and pending requests
type batchRequestsClient struct {
	address string

	conn   *grpc.ClientConn
	client batch.BatchService_StreamBatchRequestClient

	// requestID -> Request pending map
	batched sync.Map

	// requestID sequence
	idSeq uint64

	// 0 -> running
	closed int32

	clientLock sync.Mutex
}

func (c *batchRequestsClient) isClosed() bool {
	return atomic.LoadInt32(&c.closed) != 0
}

// close pending requests with err and delete requests
func (c *batchRequestsClient) failPendingRequests(err error) {
	c.batched.Range(func(k, v interface{}) bool {
		id, _ := k.(uint64)
		entry, _ := v.(*batchRequestEntry)
		entry.err = err
		close(entry.res)
		c.batched.Delete(id)
		return true
	})
}

// background receive loop, should be called with goroutine
func (c *batchRequestsClient) batchReceiveLoop() {
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Error("recover from batch receive loop", zap.Stack("stack"), zap.Any("recover", r))
			go c.batchReceiveLoop()
		}
	}()

	for {
		resp, err := c.client.Recv()
		if err != nil {
			for {
				if c.isClosed() {
					return
				}

				logger.GetLogger().Error("batch receive loop error", zap.Stack("stack"), zap.Error(err))

				c.clientLock.Lock()
				c.failPendingRequests(err)

				// grpc will handle transport layer(grpc.CliConn) re-connect
				gCli := batch.NewBatchServiceClient(c.conn)
				gStreamClient, err := gCli.StreamBatchRequest(context.TODO())
				c.clientLock.Unlock()

				if err == nil {
					logger.GetLogger().Info("re-create stream baseClient")
					c.client = gStreamClient
					break
				}

				logger.GetLogger().Error("re-create stream baseClient error", zap.Stack("stack"), zap.Error(err))

				time.Sleep(time.Second)
			}
			continue
		}

		responses := resp.GetResponses()
		for i, requestID := range resp.GetRequestIDs() {
			value, ok := c.batched.Load(requestID)
			if !ok {
				// should never happen
				panic("batch receive loop receives an unknown response")
			}

			// type assert to convert sync.Map value interface to *batchRequestEntry
			entry := value.(*batchRequestEntry)
			if !entry.isCanceled() {
				entry.res <- responses[i]
			}

			c.batched.Delete(requestID)
		}

		// todo adjust transport layer load
	}
}
