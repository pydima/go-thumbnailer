package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/streadway/amqp"
)

type rabbitMQBackend struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	queue      string
	once       sync.Once
	msgs       <-chan amqp.Delivery
	deliveries map[string]*amqp.Delivery
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connection(name string) (conn *amqp.Connection, ch *amqp.Channel) {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.Qos(config.Base.Workers, 0, false)
	failOnError(err, "Failed to set Qos")

	_, err = ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return conn, ch
}

func (mb *rabbitMQBackend) Get() (*Task, error) {
	mb.once.Do(
		func() {
			msgs, err := mb.channel.Consume(
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
		})

	t := New()
	msg := <-mb.msgs
	err := json.Unmarshal(msg.Body, t)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal data: %s", err.Error())
		return nil, err
	}

	mb.deliveries[t.TaskID] = &msg
	return t, nil
}

func (mb *rabbitMQBackend) Put(t *Task) {
	data, err := json.Marshal(*t)
	failOnError(err, "cannot marshal data")
	err = mb.channel.Publish(
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

func (mb *rabbitMQBackend) Close() {
	if mb.channel != nil {
		mb.channel.Close()
	}
	if mb.conn != nil {
		mb.conn.Close()
	}
}

func (mb *rabbitMQBackend) Complete(t *Task) {
	d, ok := mb.deliveries[t.TaskID]
	if !ok {
		return
	}
	delete(mb.deliveries, t.TaskID)
	d.Ack(false)
}
