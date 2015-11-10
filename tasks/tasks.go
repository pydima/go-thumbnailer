package tasks

import (
	"errors"
	"log"

	"github.com/pydima/go-thumbnailer/config"
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

type ImageSource struct {
	Path       string
	Identifier string
}

type Task struct {
	Images    []ImageSource
	TaskID    string
	NotifyUrl string
}

type Tasker interface {
	Get() *Task
	Put(*Task)
	Close()
}

func NewBackend(bType string) (t Tasker, err error) {
	switch bType {
	case "Memory":
		t = &MemoryBackend{make(chan *Task)}
	case "RabbitMQ":
		conn, ch, q := get_connection()
		t = &RabbitMQBackend{conn, ch, q}
	default:
		err = errors.New("Unknown backend.")
	}
	return
}
