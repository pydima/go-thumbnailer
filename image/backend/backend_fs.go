package backend

import (
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FSBackend struct {
	BasePath string
}

func (fb FSBackend) generatePath(hash string) (path string) {
	return filepath.Join(fb.BasePath, hash[:2], hash[2:4], hash[4:])
}

func (fb FSBackend) exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fb FSBackend) Save(img []byte, filename string) (string, error) {
	s := sha1.Sum(img)
	hash := base64.URLEncoding.EncodeToString(s[:])
	dir := fb.generatePath(hash)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, filename)
	if fb.exists(path) {
		return path, nil
	}

	f, err := ioutil.TempFile(dir, "tmp_")
	if err != nil {
		return "", err
	}

	defer func() {
		name := f.Name()
		if fb.exists(name) {
			os.Remove(name)
		}
	}()

	if err := f.Chmod(0644); err != nil {
		f.Close()
		return "", err
	}

	if _, err := f.Write(img); err != nil {
		f.Close()
		return "", err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	if err := os.Rename(f.Name(), path); err != nil {
		return "", err
	}

	return path, nil
}
