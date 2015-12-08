package image

import (
	"io/ioutil"
	"path/filepath"
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

func TestCheckImageFormat(t *testing.T) {
	bp := "/go/src/github.com/pydima/go-thumbnailer/testdata/"

	readAndCheckFile := func(name string) []byte {
		b, err := ioutil.ReadFile(filepath.Join(bp, name))
		if err != nil {
			t.Errorf("Cannot read the test file")
		}
		return b
	}

	b := readAndCheckFile("jpg.jpg")
	if f := checkImageFormat(b); f != JPG {
		t.Errorf("Invalid detection of jpg file.")
	}

	b = readAndCheckFile("png.png")
	if f := checkImageFormat(b); f != PNG {
		t.Errorf("Invalid detection of png file.")
	}

	b = readAndCheckFile("gif.gif")
	if f := checkImageFormat(b); f != GIF {
		t.Errorf("Invalid detection of gif file.")
	}

	b = readAndCheckFile("bmp.bmp")
	if f := checkImageFormat(b); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}

	var bt []byte
	if f := checkImageFormat(bt); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}
}
