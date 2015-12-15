package backend

type ImageBackender interface {
	Save(imgs map[string][]byte) (paths []string, err error)
}
