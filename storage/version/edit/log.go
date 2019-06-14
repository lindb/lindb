package edit

type Log interface {
	Name() string
	Encode()
	Decode()
}
type LogType string

type NewLogFunc func() Log
//
var newLogFuncMap = make(map[string]Log)

func RegisterLogType(logType string, fn Log) {
	if _, ok := newLogFuncMap[logType]; ok {
		panic("log type already registered: " + logType)
	}
	newLogFuncMap[logType] = fn
}

func NewLog() {

}
