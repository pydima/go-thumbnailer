package image

import (
	"io/ioutil"
	"path/filepath"
	"strings"
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
		t.Fatalf("Cannot read the test file")
	}
	return b
}

func TestConstructName(t *testing.T) {
	images := []string{"jpg.jpg", "png.png", "gif.gif"}
	for _, i := range images {
		want := i
		name := strings.Split(i, ".")
		if res := constructName(name[0], name[1]); res != want {
			t.Errorf("constructName(); want - %s, get - %s", want, res)
		}
	}
}

func TestImageFormat(t *testing.T) {
	b := readAndCheckFile("jpg.jpg", t)
	if f := imageFormat(b); f != JPG {
		t.Errorf("Invalid detection of jpg file.")
	}

	b = readAndCheckFile("png.png", t)
	if f := imageFormat(b); f != PNG {
		t.Errorf("Invalid detection of png file.")
	}

	b = readAndCheckFile("gif.gif", t)
	if f := imageFormat(b); f != GIF {
		t.Errorf("Invalid detection of gif file.")
	}

	b = readAndCheckFile("bmp.bmp", t)
	if f := imageFormat(b); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}

	var bt []byte
	if f := imageFormat(bt); f != UNKNOWN {
		t.Errorf("Invalid detection of unknown file format.")
	}
}

func checkDimensions(b []byte, width, height int, t *testing.T, exact bool) {
	w, h, err := imageDimensions(b)
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
	checkDimensions(b, 1634, 2224, t, true)

	b = readAndCheckFile("png.png", t)
	checkDimensions(b, 1634, 2224, t, true)

	b = readAndCheckFile("gif.gif", t)
	checkDimensions(b, 1634, 2224, t, true)

	b = readAndCheckFile("bmp.bmp", t)
	_, _, err := imageDimensions(b)
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
	if f := imageFormat(res); f != PNG {
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
		Background: bimg.Color{R: 255, G: 255, B: 255},
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
		Background: bimg.Color{R: 255, G: 255, B: 255},
		Type:       3,
	}
	b := readAndCheckFile("jpg.jpg", t)
	res, err := bimg.Resize(b, options)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if f := imageFormat(res); f != PNG {
		t.Errorf("Invalid image format.")
	}
}

func TestProcessImage(t *testing.T) {
	options := bimg.Options{
		Width:      100,
		Height:     100,
		Enlarge:    true,
		Quality:    95,
		Background: bimg.Color{R: 255, G: 255, B: 255},
		Type:       3,
	}
	b := readAndCheckFile("gif.gif", t)
	ip := config.ImageParam{Name: "test", Extension: "png"}
	ip.Options = options
	res, err := processImage(b, ip)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if f := imageFormat(res); f != PNG {
		t.Errorf("Invalid image format.")
	}
	checkDimensions(res, 100, 100, t, false)
}

func TestCreateThumbnails(t *testing.T) {
	b := readAndCheckFile("png.png", t)
	res, err := CreateThumbnails(b)
	if err != nil {
		t.Errorf("Got error: %s", err)
	}

	if len(res) < 2 {
		t.Errorf("Got not enough thumbnails.")
	}

	for k, img := range res {
		var it imageType
		if strings.HasPrefix(k, "original") {
			it = PNG
		} else {
			it = JPG
		}
		if f := imageFormat(img); f != it {
			t.Errorf("Invalid image format, got %d expected %d.", f, it)
		}
	}
}
