package image

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FSBackend struct {
	BasePath string
}

func (fb FSBackend) createPath(hash string) (path string) {
	return filepath.Join(fb.BasePath, hash[:2], hash[2:4], hash[4:])
}

func (fb FSBackend) Exists(hash string) bool {
	path := fb.createPath(hash)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fb FSBackend) Save(img []byte, hash string) (path string, err error) {
	path = fb.createPath(hash)
	if fb.Exists(hash) {
		return path, nil
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	if err = ioutil.WriteFile(path, img, 0644); err != nil {
		return "", err
	}

	return path, nil
}
