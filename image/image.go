package image

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/daddye/vips"

	"github.com/pydima/go-thumbnailer/config"
)

var (
	MARKER_JPG = []byte{0xff, 0xd8}
	MARKER_PNG = []byte{0x89, 0x50}
	MARKER_GIF = []byte{0x47, 0x49}
)

type ImageType int

const (
	UNKNOWN ImageType = iota
	JPG
	PNG
	GIF
)

type Image struct {
	width  uint
	height uint
	path   string
}

type InvalidImage struct {
	err string
}

func (e InvalidImage) Error() string {
	return e.err
}

func checkExtension(n string) error {
	for _, ext := range config.Base.ValidExtensions {
		if strings.HasSuffix(strings.ToLower(n), ext) {
			return nil
		}
	}
	return InvalidImage{fmt.Sprintf("Extension %s is not supported.", filepath.Ext(n))}
}

func checkImageFormat(img []byte) ImageType {
	if len(img) < 2 {
		return UNKNOWN
	}

	switch {
	case bytes.Equal(img[:2], MARKER_JPG):
		return JPG
	case bytes.Equal(img[:2], MARKER_PNG):
		return PNG
	case bytes.Equal(img[:2], MARKER_GIF):
		return GIF
	default:
		return UNKNOWN
	}
}

func getImageDimensions(img []byte) (width, height int, err error) {
	r := bytes.NewReader(img)
	conf, _, err := image.DecodeConfig(r)
	return conf.Width, conf.Height, err
}

// vips doesn't support gif natively, so have to convert it with slow standart library
func convertGifToPng(img []byte) ([]byte, error) {
	r := bytes.NewReader(img)
	i, err := gif.Decode(r)
	if err != nil {
		return nil, err
	}

	res := new(bytes.Buffer)
	err = png.Encode(res, i)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func ProcessImage(i io.ReadCloser) (path string, err error) {
	options := vips.Options{
		Width:        100,
		Height:       100,
		Crop:         true,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	path = "res.jpg"

	f, _ := os.Create(path)
	w := bufio.NewWriter(f)
	input, _ := ioutil.ReadAll(i)

	buf, err := vips.Resize(input, options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	w.Write(buf)
	w.Flush()
	fmt.Println("Done!")
	return
}
