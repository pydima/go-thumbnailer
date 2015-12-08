package backend

type ImageBackender interface {
	Exists(hash string) bool
	Save(img []byte, hash string) (path string, err error)
	Load(hash string) (img []byte)
}
