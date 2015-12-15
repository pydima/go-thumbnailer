package backend

type ImageBackender interface {
	Save(img []byte) (path string, err error)
}
