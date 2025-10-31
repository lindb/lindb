package input

import (
	"sync"
)

type InputManager interface {
	GetInputHandler(stream string) InputHandler
}

type inputManager struct {
	inputHandlers map[string]InputHandler

	mutex sync.Mutex
}

func NewInputManager() InputManager {
	return &inputManager{
		inputHandlers: make(map[string]InputHandler),
	}
}

func (mgr *inputManager) GetInputHandler(stream string) InputHandler {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	handle, ok := mgr.inputHandlers[stream]
	if ok {
		return handle
	}
	handle = NewInputHandler()
	mgr.inputHandlers[stream] = handle
	return handle
}
