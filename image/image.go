package image

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/daddye/vips"

	"github.com/pydima/go-thumbnailer/config"
)

type Image struct {
	width  uint
	height uint
	path   string
}

type InvalidExtension struct {
	err string
}

func (e InvalidExtension) Error() string {
	return e.err
}

func checkExtension(n string) error {
	for _, ext := range config.Base.ValidExtensions {
		if strings.HasSuffix(strings.ToLower(n), ext) {
			return nil
		}
	}
	return InvalidExtension{fmt.Sprintf("Cannot handle %s extension")}
}

func ProcessImage(i io.ReadCloser) (path string, err error) {
	options := vips.Options{
		Width:        100,
		Height:       100,
		Crop:         false,
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
