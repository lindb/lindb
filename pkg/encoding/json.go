package encoding

import (
	"encoding/json"

	"github.com/lindb/lindb/pkg/logger"
)

var log = logger.GetLogger("encoding")

// JSONMarshal returns the JSON encoding of v.
func JSONMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("json marshal error")
	}
	return data
}

// JSONUnmarshal parses the JSON-encoded data and stores the result
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
