package tasks

import (
	"testing"
)

func TestPutGetMemory(t *testing.T) {
	task := New()

	b, err := NewBackend("Memory")
	if err != nil {
		t.Error(err)
	}

	go b.Put(task)
	task2, _ := b.Get()
	if task.TaskID != task2.TaskID {
		t.Error("Tasks are not the same.")
	}
}
