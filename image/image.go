package image

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	return InvalidImage{fmt.Sprintf("Cannot handle %s extension")}
}

func checkImageFormat(b []byte) ImageType {
	if len(b) < 2 {
		return UNKNOWN
	}

	switch {
	case bytes.Equal(b[:2], MARKER_JPG):
		return JPG
	case bytes.Equal(b[:2], MARKER_PNG):
		return PNG
	case bytes.Equal(b[:2], MARKER_GIF):
		return GIF
	default:
		return UNKNOWN
	}
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
