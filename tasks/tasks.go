package tasks

import (
	"log"

	"github.com/pydima/go-thumbnailer/config"
)

var Backend Tasker

func init() {
	bt := config.Base.TaskBackend

	switch bt {
	case "Memory":
		Backend = &MemoryBackend{make(chan *Task)}
	case "RabbitMQ":
		conn, ch, q := get_connection()
		Backend = &RabbitMQBackend{conn, ch, q}
	default:
		log.Fatal("Unknown backend.")
	}
}

type ImageSource struct {
	Path       string
	Identifier string
}

type Task struct {
	Images []ImageSource
	TaskID string
}

type Tasker interface {
	Get() *Task
	Put(*Task)
	Close()
}
