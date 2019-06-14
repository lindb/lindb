package broker

type Config struct {
	HTTP HTTP
}

type HTTP struct {
	Port   int32
	Static string
}
