package tasks

import (
	"testing"

	"github.com/streadway/amqp"
)

var (
	rabbitBackend *rabbitMQBackend
	queue         = "test_images_queue"
)

func init() {
	conn, ch := connection(queue)
	rabbitBackend = &rabbitMQBackend{conn: conn, channel: ch, queue: queue, deliveries: make(map[string]*amqp.Delivery)}
}

func TestPutGetRabbitMQ(t *testing.T) {
	defer rabbitBackend.channel.QueuePurge(queue, true)

	task := New()

	go rabbitBackend.Put(task)

	task2, err := rabbitBackend.Get()
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}

	rabbitBackend.Complete(task2)

	if task.TaskID != task2.TaskID {
		t.Errorf("Tasks are not the same. (%s -> %s)", task.TaskID, task2.TaskID)
	}
}

func TestRabbitAckLate(t *testing.T) {
	defer rabbitBackend.channel.QueuePurge(queue, true)

	task := New()

	go rabbitBackend.Put(task)
	task2, err := rabbitBackend.Get()
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}

	_, ok := rabbitBackend.deliveries[task.TaskID]
	if !ok {
		t.Errorf("Couldn't find row in deliveries.")
	}

	rabbitBackend.Complete(task2)

	_, ok = rabbitBackend.deliveries[task.TaskID]
	if ok {
		t.Errorf("Found row in deliveries.")
	}

}
