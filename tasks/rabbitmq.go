package tasks

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQBackend struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	msgs       <-chan amqp.Delivery
	deliveries map[string]*amqp.Delivery
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func connection(name string) (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return conn, ch, &q
}

func (mb *RabbitMQBackend) Get() *Task {
	if mb.msgs == nil {
		msgs, err := mb.channel.Consume(
			mb.queue.Name, // queue
			"",            // consumer
			false,         // auto-ack
			false,         // exclusive
			false,         // no-local
			false,         // no-wait
			nil,           // args
		)
		failOnError(err, "Failed to register a consumer")
		mb.msgs = msgs
	}

	t := New()
	msg := <-mb.msgs
	err := json.Unmarshal(msg.Body, t)
	failOnError(err, "Failed to unmarshal data")

	mb.deliveries[t.TaskID] = &msg
	return t
}

func (mb *RabbitMQBackend) Put(t *Task) {
	data, err := json.Marshal(*t)
	failOnError(err, "cannot marshal data")
	err = mb.channel.Publish(
		"",            // exchange
		mb.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")
	return
}

func (mb *RabbitMQBackend) Close() {
	mb.channel.Close()
	mb.conn.Close()
}

func (mb *RabbitMQBackend) Complete(t *Task) {
	d, ok := mb.deliveries[t.TaskID]
	if !ok {
		return
	}
	delete(mb.deliveries, t.TaskID)
	d.Ack(false)
}
