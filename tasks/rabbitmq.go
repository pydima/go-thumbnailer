package tasks

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQBackend struct {
	conn       *amqp.Connection
	pubChannel *amqp.Channel
	subChannel *amqp.Channel
	queue      string
	msgs       <-chan amqp.Delivery
	deliveries map[string]*amqp.Delivery
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func connection(name string) (conn *amqp.Connection, pubCh, subCh *amqp.Channel) {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	pubCh, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	subCh, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	_, err = pubCh.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return conn, pubCh, subCh
}

func (mb *RabbitMQBackend) Get() *Task {
	if mb.msgs == nil {
		msgs, err := mb.subChannel.Consume(
			mb.queue, // queue
			"",       // consumer
			false,    // auto-ack
			false,    // exclusive
			false,    // no-local
			false,    // no-wait
			nil,      // args
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
	err = mb.pubChannel.Publish(
		"",       // exchange
		mb.queue, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
	failOnError(err, "Failed to publish a message")
	return
}

func (mb *RabbitMQBackend) Close() {
	mb.pubChannel.Close()
	mb.subChannel.Close()
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
