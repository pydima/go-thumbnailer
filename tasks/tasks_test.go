package tasks

import "testing"

func TestGenerateTaskID(t *testing.T) {
	task := New()
	task2 := New()
	if task.TaskID == task2.TaskID {
		t.Errorf("Ids are the same")
	}
}
