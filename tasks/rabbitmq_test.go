package tasks

import (
	"testing"
)

var RabbitBackend *RabbitMQBackend

func init() {
	b, err := NewBackend("RabbitMQ")
	if err != nil {
		panic(err)
	}
	RabbitBackend = b.(*RabbitMQBackend)
}

func TestPutGetRabbitMQ(t *testing.T) {
	task := &Task{TaskID: "test_task"}

	go RabbitBackend.Put(task)

	task2 := RabbitBackend.Get()
	if task.TaskID != task2.TaskID {
		t.Error("Tasks are not the same.")
	}
}

func TestRabbitAckLate(t *testing.T) {
	task := &Task{TaskID: "test_task"}

	go RabbitBackend.Put(task)
	task2 := RabbitBackend.Get()

	_, ok := RabbitBackend.deliveries[task.TaskID]
	if !ok {
		t.Errorf("Couldn't find row in deliveries.")
	}

	RabbitBackend.Complete(task2)

	_, ok = RabbitBackend.deliveries[task.TaskID]
	if ok {
		t.Errorf("Found row in deliveries.")
	}

}
