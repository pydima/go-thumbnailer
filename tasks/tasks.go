package tasks

import (
	"errors"
	"log"

	"github.com/streadway/amqp"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/utils"
)

// Backend is a singleton with current implementation
// of Tasker interface, implemenation is specified in the config
var Backend Tasker

func init() {
	var err error
	bt := config.Base.TaskBackend
	Backend, err = newBackend(bt)
	if err != nil {
		log.Fatal(err)
	}
}

// New returns pointer to a new initialized task
func New() *Task {
	return &Task{
		Images: make([]ImageSource, 3),
		TaskID: utils.UUID(),
	}
}

// ImageSource contains url with a source image and image's identifier
// which might be used to download image from the url which we already
// processed. For example, if the image from url changes and we want to
// update data, it's possible to just specify the same url with a
// different identifier.
type ImageSource struct {
	Path       string
	Identifier string
}

// Task is a general structure that contains all information about
// the image which we need to process.
type Task struct {
	Images    []ImageSource
	TaskID    string
	ID        string
	NotifyURL string
}

// Tasker defines an interface for task backends.
// It is possible to use such task backends as
// Memory, RabbitMQ, Redis, etc.
type Tasker interface {
	Get() (*Task, error)
	Put(*Task)
	Close()
	Complete(*Task)
}

func newBackend(bType string) (t Tasker, err error) {
	switch bType {
	case "Memory":
		t = &memoryBackend{make(chan *Task, 100)}
	case "RabbitMQ":
		queue := "images"
		conn, ch := connection(queue)
		t = &rabbitMQBackend{conn: conn, channel: ch, queue: queue, deliveries: make(map[string]*amqp.Delivery)}
	default:
		err = errors.New("Unknown backend.")
	}
	return
}
