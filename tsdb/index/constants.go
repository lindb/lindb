package index

import (
	"encoding/json"
	"math"
)

const (
	NotFoundTagsID uint32 = math.MaxInt32

	NotFoundMetricID uint32 = math.MaxUint32

	MetricSequenceIDKey = 0

	NotFoundFieldID uint32 = math.MaxUint32
)

func StringToMap(tags string) map[string]string {
	var tagsMap map[string]string
	_ = json.Unmarshal([]byte(tags), &tagsMap)
	return tagsMap
}

func MapToString(tagsMap map[string]string) string {
	b, _ := json.Marshal(tagsMap)
	return string(b)
}
