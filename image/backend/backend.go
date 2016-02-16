package backend

import (
	"errors"
	"log"
	"time"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/utils"
)

// ImageBackender defines interface which should be implemented by
// all storage backends
type ImageBackender interface {
	Save(imgs map[string][]byte) (paths []string, err error)
}

// AlreadyExistsError will be returned in
// case the image already exists
type AlreadyExistsError struct {
	err  error
	Path string
}

func (e *AlreadyExistsError) Error() string {
	return e.err.Error()
}

// ImageBackend is a singleton with the current
// ImageBackender implementation
var ImageBackend ImageBackender

func init() {
	var err error
	bt := config.Base.ImageBackend
	ImageBackend, err = newBackend(bt)
	if err != nil {
		log.Fatal(err)
	}
}

func newBackend(bType string) (t ImageBackender, err error) {
	switch bType {
	case "FS":
		t = backendFS{
			BasePath: config.Base.MediaRoot,
			TmpDir:   config.Base.TmpDir,
		}
		ticker := time.NewTicker(10 * time.Minute)
		go func() {
			for {
				select {
				case <-ticker.C:
					imageGC(config.Base.TmpDir)
				case <-utils.STOP:
					ticker.Stop()
					return
				}
			}
		}()
	default:
		err = errors.New("Unknown backend.")
	}
	return
}
