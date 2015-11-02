package tasks

import (
	"github.com/pydima/go-thumbnailer/config"
)

var Backend Tasker

func init() {
	bt := config.Base.TaskBackend

	if bt == "Memory" {
		Backend = &MemoryBackend{make(chan *Task)}
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
}
