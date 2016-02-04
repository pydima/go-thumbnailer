package tasks

type MemoryBackend struct {
	tasks chan *Task
}

func (mb *MemoryBackend) Get() (*Task, error) {
	res := <-mb.tasks
	return res, nil
}

func (mb *MemoryBackend) Put(t *Task) {
	mb.tasks <- t
	return
}

func (mb *MemoryBackend) Close() {}

func (mb *MemoryBackend) Complete(t *Task) {}
