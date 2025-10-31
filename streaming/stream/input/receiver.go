package input

type Receiver interface {
	Receive(event any)
}
