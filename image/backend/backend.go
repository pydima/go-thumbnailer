package backend

import (
	"errors"
	"log"

	"github.com/pydima/go-thumbnailer/config"
)

type ImageBackender interface {
	Save(imgs map[string][]byte) (paths []string, err error)
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
	default:
		err = errors.New("Unknown backend.")
	}
	return
}
