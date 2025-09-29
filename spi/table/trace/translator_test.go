package trace

import (
	"encoding/hex"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestTraceID(t *testing.T) {
	var traceID trace.TraceID
	b, err := hex.DecodeString("0fa9e36f12d744ff243a73dde06add94")
	if err != nil {
		panic(err)
	}
	if len(b) != len(traceID) {
		return
	}
	copy(traceID[:], b)
	fmt.Println(traceID.String())
}
