package monitoring

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPCollector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = ioutil.ReadAll(r.Body)
		}))
	defer ts.Close()

	collector := NewHTTPCollector(
		ctx,
		ts.URL,
		ts.URL,
		time.Millisecond*100,
	)
	defer collector.Stop()

	go collector.Run()

	time.Sleep(200 * time.Millisecond)
}

func TestHTTPCollector_http_fail(t *testing.T) {
	defer func() {
		doGet = http.Get
		readBody = ioutil.ReadAll
		newRequest = http.NewRequest
		doRequest = http.DefaultClient.Do
	}()

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = ioutil.ReadAll(r.Body)
		}))
	defer ts.Close()

	collector := NewHTTPCollector(
		context.TODO(),
		ts.URL,
		ts.URL,
		time.Millisecond*100,
	)

	c := collector.(*httpCollector)
	// case 1: get err
	doGet = func(url string) (resp *http.Response, err error) {
		return nil, fmt.Errorf("err")
	}
	c.collect()
	doGet = http.Get
	// case 2: read body err
	readBody = func(r io.Reader) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	c.collect()
	readBody = ioutil.ReadAll
	// case 3: new request err
	newRequest = func(method, url string, body io.Reader) (request *http.Request, err error) {
		return nil, fmt.Errorf("err")
	}
	c.collect()
	newRequest = http.NewRequest
	// case 4: do request err
	doRequest = func(req *http.Request) (response *http.Response, err error) {
		return nil, fmt.Errorf("err")
	}
	c.collect()
}
