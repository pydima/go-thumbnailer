package image

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/h2non/bimg"

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

func checkDimensions(b []byte, width, height int, t *testing.T, exact bool) {
	w, h, err := getImageDimensions(b)
	if err != nil {
		t.Errorf("Cannot get dimensions.")
		return
	}

	var invalid bool
	if exact {
		if w != width || h != height {
			invalid = true
		}
	} else {
		if w != width && h != height {
			invalid = true
		}
	}

	if invalid {
		t.Errorf("Got invalid dimensions, width: %d, height: %d", w, h)
	}
}

func TestGetImageDimensions(t *testing.T) {
	b := readAndCheckFile("jpg.jpg", t)
	checkDimensions(b, 1431, 901, t, true)

	b = readAndCheckFile("png.png", t)
	checkDimensions(b, 1634, 2224, t, true)

	b = readAndCheckFile("gif.gif", t)
	checkDimensions(b, 450, 159, t, true)

	b = readAndCheckFile("bmp.bmp", t)
	_, _, err := getImageDimensions(b)
	if err == nil {
		t.Errorf("Should return an error, since cannot get dimensions.")
	}
}

func TestConvertGifToPng(t *testing.T) {
	b := readAndCheckFile("gif.gif", t)
	res, err := convertGifToPng(b)
	if err != nil {
		t.Errorf("Cannot convert image to png.")
	}
	if f := checkImageFormat(res); f != PNG {
		t.Errorf("Convertation failed.")
	}

	b = readAndCheckFile("png.png", t)
	res, err = convertGifToPng(b)
	if err == nil {
		t.Errorf("Successfully converted png, but should support only gif.")
	}
}

func TestCreateThumbnail(t *testing.T) {
	options := bimg.Options{
		Width:      100,
		Height:     100,
		Enlarge:    true,
		Quality:    95,
		Background: bimg.Color{255, 255, 255},
	}

	b := readAndCheckFile("png.png", t)
	res, err := bimg.Resize(b, options)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	checkDimensions(res, 100, 100, t, false)
}

func TestConvertJpgToPng(t *testing.T) {
	options := bimg.Options{
		Width:      100,
		Height:     100,
		Enlarge:    true,
		Quality:    95,
		Background: bimg.Color{255, 255, 255},
		Type:       3,
	}
	b := readAndCheckFile("jpg.jpg", t)
	res, err := bimg.Resize(b, options)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if f := checkImageFormat(res); f != PNG {
		t.Errorf("Invalid image format.")
	}
}
