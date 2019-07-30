package parallel

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestTaskManager_ClientStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentNode := models.Node{IP: "1.1.1.1", Port: 8000}
	taskSender := NewMockTaskSenderManager(ctrl)

	taskManager1 := NewTaskManager(currentNode, taskSender)

	taskCtx := newTaskContext("xxx", "parentTaskID", "parentNode", 2)
	taskManager1.Submit(taskCtx)

	assert.Equal(t, taskCtx, taskManager1.Get("xxx"))
	assert.Nil(t, taskManager1.Get("xxx11"))
	assert.Equal(t, taskSender, taskManager1.GetTaskSenderManager())

	taskManager2 := taskManager1.(*taskManager)
	taskManager2.tasks.Store("xxx11", nil)
	assert.Nil(t, taskManager1.Get("xxx11"))

	taskCtx = newTaskContext("taskID", "parentTaskID", "parentNode", 2)
	taskManager1.Submit(taskCtx)
	assert.Equal(t, taskCtx, taskManager1.Get("taskID"))
	taskManager1.Complete("taskID")
	assert.Nil(t, taskManager1.Get("taskID"))

	assert.Equal(t, "1.1.1.1:8000-1", taskManager1.AllocTaskID())
	assert.Equal(t, "1.1.1.1:8000-2", taskManager1.AllocTaskID())
	assert.Equal(t, "1.1.1.1:8000-3", taskManager1.AllocTaskID())
}
