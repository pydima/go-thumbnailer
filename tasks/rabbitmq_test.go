package tasks

import (
	"testing"

	"github.com/streadway/amqp"
)

var (
	RabbitBackend *RabbitMQBackend
	queue         string = "test_images_queue"
)

func init() {
	conn, ch := connection(queue)
	RabbitBackend = &RabbitMQBackend{conn: conn, channel: ch, queue: queue, deliveries: make(map[string]*amqp.Delivery)}
}

func TestPutGetRabbitMQ(t *testing.T) {
	defer RabbitBackend.channel.QueuePurge(queue, true)

	task := New()

	go RabbitBackend.Put(task)

	task2, err := RabbitBackend.Get()
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}

	RabbitBackend.Complete(task2)

	if task.TaskID != task2.TaskID {
		t.Errorf("Tasks are not the same. (%s -> %s)", task.TaskID, task2.TaskID)
	}
}

func TestRabbitAckLate(t *testing.T) {
	defer RabbitBackend.channel.QueuePurge(queue, true)

	task := New()

	go RabbitBackend.Put(task)
	task2, err := RabbitBackend.Get()
	if err != nil {
		t.Errorf("Got error: %s", err.Error())
		return
	}

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
