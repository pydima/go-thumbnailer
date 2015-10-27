package image

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/daddye/vips"
)

type Image struct{}

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
