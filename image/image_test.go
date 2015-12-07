package image

import (
	"testing"

	"github.com/pydima/go-thumbnailer/config"
)

func TestCheckExtension(t *testing.T) {
	config.Base.ValidExtensions = []string{"jpg"}
	err := checkExtension("file.jpg")
	if err != nil {
		t.Errorf("Cannot detect valid extension.")
	}

	err = checkExtension("file.svg")
	if err == nil {
		t.Errorf("Detect invalid extension.")
	}

	err = checkExtension("file.JPG")
	if err != nil {
		t.Errorf("Cannot detect extension in uppercase.")
	}
}
