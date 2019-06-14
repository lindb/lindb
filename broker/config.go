package broker

type Config struct {
	Http Http
}

type Http struct {
	Port int32
	Static string
}
