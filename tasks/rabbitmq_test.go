package tasks

import (
	"testing"
)

func TestPutGetRabbitMQ(t *testing.T) {
	task := &Task{TaskID: "test_task"}

	b, err := NewBackend("RabbitMQ")
	if err != nil {
		t.Error(err)
	}

	go b.Put(task)
	task2 := b.Get()
	if task.TaskID != task2.TaskID {
		t.Error("Tasks are not the same.")
	}
}
