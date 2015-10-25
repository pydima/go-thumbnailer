package image

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Image struct{}

func ProcessImage(i io.ReadCloser) (Image, err error) {
	f, _ := os.Create("/home/home/res.jpg")
	w := bufio.NewWriter(f)
	io.Copy(w, i)
	w.Flush()
	fmt.Println("Done!")
	return
}
