package tasks

type MemoryBackend struct {
	tasks chan *Task
}

func (mb *MemoryBackend) Get() *Task {
	res := <-mb.tasks
	return res
}

func (mb *MemoryBackend) Put(t *Task) {
	mb.tasks <- t
	return
}
