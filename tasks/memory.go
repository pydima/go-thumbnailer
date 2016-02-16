package tasks

type memoryBackend struct {
	tasks chan *Task
}

func (mb *memoryBackend) Get() (*Task, error) {
	res := <-mb.tasks
	return res, nil
}

func (mb *memoryBackend) Put(t *Task) {
	mb.tasks <- t
	return
}

func (mb *memoryBackend) Close() {}

func (mb *memoryBackend) Complete(t *Task) {}
