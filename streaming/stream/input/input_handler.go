package input

type InputHandler interface {
	Send(event any)
	Subscribe(receiver Receiver)
}

type inputHandler struct {
	receivers []Receiver
}

func NewInputHandler() InputHandler {
	return &inputHandler{}
}

func (h *inputHandler) Subscribe(receiver Receiver) {
	h.receivers = append(h.receivers, receiver)
}

func (h *inputHandler) Send(event any) {
	// fmt.Printf("current:%v, send event:%v\n", h, event)
	for _, r := range h.receivers {
		r.Receive(event)
	}
}
