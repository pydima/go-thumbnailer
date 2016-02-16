package backend

import (
	"errors"
	"log"
	"time"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/utils"
)

type ImageBackender interface {
	Save(imgs map[string][]byte) (paths []string, err error)
}

type AlreadyExistsError struct {
	err  error
	Path string
}

func (e *AlreadyExistsError) Error() string {
	return e.err.Error()
}

var ImageBackend ImageBackender

func init() {
	var err error
	bt := config.Base.ImageBackend
	ImageBackend, err = NewBackend(bt)
	if err != nil {
		log.Fatal(err)
	}
}

func NewBackend(bType string) (t ImageBackender, err error) {
	switch bType {
	case "FS":
		t = FSBackend{
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
