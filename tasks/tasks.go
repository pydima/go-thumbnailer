package tasks

import (
	"errors"
	"log"

	"github.com/streadway/amqp"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/utils"
)

var Backend Tasker

func init() {
	var err error
	bt := config.Base.TaskBackend
	Backend, err = NewBackend(bt)
	if err != nil {
		log.Fatal(err)
	}
}

func New() *Task {
	return &Task{
		Images: make([]ImageSource, 3),
		TaskID: utils.UUID(),
	}
}

type ImageSource struct {
	Path       string
	Identifier string
}

type Task struct {
	Images    []ImageSource
	TaskID    string
	ID        string
	NotifyUrl string
}

type Tasker interface {
	Get() (*Task, error)
	Put(*Task)
	Close()
	Complete(*Task)
}

func NewBackend(bType string) (t Tasker, err error) {
	switch bType {
	case "Memory":
		t = &MemoryBackend{make(chan *Task, 100)}
	case "RabbitMQ":
		queue := "images"
		conn, ch := connection(queue)
		t = &RabbitMQBackend{conn: conn, channel: ch, queue: queue, deliveries: make(map[string]*amqp.Delivery)}
	default:
		err = errors.New("Unknown backend.")
	}
	return
}
