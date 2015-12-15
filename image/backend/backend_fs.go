package backend

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type FSBackend struct {
	BasePath string
}

func (fb FSBackend) generatePath(hash string) (path string) {
	return filepath.Join(fb.BasePath, hash[:2], hash[2:4], hash[4:6], hash[6:])
}

func (fb FSBackend) exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fb FSBackend) generateDest(imgs map[string][]byte, basePath string) []string {
	res := make([]string, 0, len(imgs))
	for k, _ := range imgs {
		res = append(res, filepath.Join(basePath, k))
	}
	return res
}

func (fb FSBackend) Save(imgs map[string][]byte) ([]string, error) {
	var (
		found bool
		img   []byte
	)
	for k, v := range imgs {
		if strings.HasPrefix(k, "original") { // image might have different extensions e.g. 'jpg', 'png'
			img = v
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("cannot find original image.")
	}

	s := sha1.Sum(img)
	hash := base64.URLEncoding.EncodeToString(s[:])
	dir := fb.generatePath(hash)
	if fb.exists(dir) {
		return fb.generateDest(imgs, dir), nil
	}

	baseTempDir := filepath.Join(fb.BasePath, "tmp")
	tmpDir, err := ioutil.TempDir(baseTempDir, time.Now().Format(time.RFC3339)+"_")
	if err != nil {
		pathErr, ok := err.(*os.PathError)
		if !ok {
			return nil, err
		}
		if pathErr.Err != syscall.ENOENT {
			return nil, err
		} else {
			if err := os.MkdirAll(baseTempDir, 0755); err != nil {
				return nil, err
			}
			tmpDir, err = ioutil.TempDir(baseTempDir, time.Now().Format(time.RFC3339)+"_")
			if err != nil {
				return nil, err
			}
		}
	}

	defer func() {
		if fb.exists(tmpDir) {
			os.RemoveAll(tmpDir)
		}
	}()

	for k, v := range imgs {
		if err := ioutil.WriteFile(filepath.Join(tmpDir, k), v, 0644); err != nil {
			return nil, err
		}
	}

	parentDir := filepath.Dir(dir)
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		return nil, err
	}

	if err := os.Rename(tmpDir, dir); err != nil {
		linkErr, ok := err.(*os.LinkError)
		if !ok {
			return nil, err
		}
		if linkErr.Err != syscall.ENOTEMPTY {
			return nil, err
		}
	}

	return fb.generateDest(imgs, dir), nil
}
