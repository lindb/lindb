package input

import (
	"sync"
)

var (
	instance *Manager
	once     sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			apps: make(map[string]InputManager),
		}
	})
	return instance
}

type Manager struct {
	apps map[string]InputManager

	mutex sync.Mutex
}

func (mgr *Manager) GetInputHandler(app, stream string) InputHandler {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	inputMgr, ok := mgr.apps[app]
	if !ok {
		inputMgr = NewInputManager()
		mgr.apps[app] = inputMgr
	}
	return inputMgr.GetInputHandler(stream)
}
