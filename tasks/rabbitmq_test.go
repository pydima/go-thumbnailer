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
	task := New()

	go RabbitBackend.Put(task)

	task2 := RabbitBackend.Get()

	if task.TaskID != task2.TaskID {
		t.Errorf("Tasks are not the same. (%s -> %s)", task.TaskID, task2.TaskID)
	}
}

func TestRabbitAckLate(t *testing.T) {
	task := New()

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
