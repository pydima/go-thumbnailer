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

func readAndCheckFile(name string, t *testing.T) []byte {
	bp := "/go/src/github.com/pydima/go-thumbnailer/testdata/"
	b, err := ioutil.ReadFile(filepath.Join(bp, name))
	if err != nil {
		t.Errorf("Cannot read the test file")
	}
	return b
}

func TestCheckImageFormat(t *testing.T) {
	b := readAndCheckFile("jpg.jpg", t)
	if f := checkImageFormat(b); f != JPG {
		t.Errorf("Invalid detection of jpg file.")
	}

	b = readAndCheckFile("png.png", t)
	if f := checkImageFormat(b); f != PNG {
		t.Errorf("Invalid detection of png file.")
	}

	b = readAndCheckFile("gif.gif", t)
	if f := checkImageFormat(b); f != GIF {
		t.Errorf("Invalid detection of gif file.")
	}

	b = readAndCheckFile("bmp.bmp", t)
	if f := checkImageFormat(b); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}

	var bt []byte
	if f := checkImageFormat(bt); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}
}

func checkDimensions(b []byte, width, height int, t *testing.T) {
	w, h, err := getImageDimensions(b)
	if err != nil {
		t.Errorf("Cannot get dimensions.")
	}
	if w != width || h != height {
		t.Errorf("Got invalid dimensions, width: %d, height: %d", width, height)
	}
}

func TestGetImageDimensions(t *testing.T) {
	b := readAndCheckFile("jpg.jpg", t)
	checkDimensions(b, 1431, 901, t)

	b = readAndCheckFile("png.png", t)
	checkDimensions(b, 1634, 2224, t)

	b = readAndCheckFile("gif.gif", t)
	checkDimensions(b, 450, 159, t)

	b = readAndCheckFile("bmp.bmp", t)
	_, _, err := getImageDimensions(b)
	if err == nil {
		t.Errorf("Should return an error, since cannot get dimensions.")
	}
}
