package monitoring

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func Test_NewRuntimeCollector(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = ioutil.ReadAll(r.Body)
		}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runtimeCollector := NewRunTimeCollector(
		ctx,
		ts.URL,
		time.Millisecond*100,
		nil)
	// manually trigger gc to cover ReportCounter
	runtime.GC()

	go runtimeCollector.Run()

	time.Sleep(time.Second * 2)
}
