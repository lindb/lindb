package execution

import "sync"

type RequestManager interface {
	CompleteRequet(id RequestID, err error)
}

type requestManager struct {
	requests map[RequestID]any
	lock     sync.RWMutex
}

func NewRequestManager() RequestManager {
	return &requestManager{
		requests: make(map[RequestID]any),
	}
}

func (mgr *requestManager) CompleteRequet(id RequestID, err error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	delete(mgr.requests, id)
}
