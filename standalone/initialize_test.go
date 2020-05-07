package standalone

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
)

func TestInitialize(t *testing.T) {
	defer func() {
		newRequest = http.NewRequest
		doRequest = http.DefaultClient.Do
	}()
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = ioutil.ReadAll(r.Body)
		}))
	defer ts.Close()

	newRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, fmt.Errorf("err")
	}
	init := newInitialize(ts.URL)

	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})
	newRequest = http.NewRequest

	doRequest = func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("err")
	}
	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})

	doRequest = http.DefaultClient.Do
	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})
}
