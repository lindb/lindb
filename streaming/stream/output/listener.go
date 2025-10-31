package output

type Listener interface {
	Receive(event any)
}
