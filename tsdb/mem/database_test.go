package mem

import (
	"testing"
	"fmt"
)

func Test_Write(t *testing.T) {
	var db = NewMemDatabase()
	f1 := db.GetTimeSeriesStore("cpu", fmt.Sprintf("host:disk:host+disk"))
	f2 := db.GetTimeSeriesStore("cpu", fmt.Sprintf("host:disk:host+disk"))
	if f1 != f2 {
		t.Errorf("get diff field store")
		return
	}
}
