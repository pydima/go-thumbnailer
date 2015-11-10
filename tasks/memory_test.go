package tasks

import (
	"testing"
)

func TestPutGet(t *testing.T) {
	var task *Task

	b, err := NewBackend("Memory")
	if err != nil {
		t.Error(err)
	}

	go b.Put(task)
	task2 := b.Get()
	if task != task2 {
		t.Error("Get different task.")
	}
}
