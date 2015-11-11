package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pydima/go-thumbnailer/image"
)

var Random *os.File

func init() {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		log.Fatal(err)
	}
	Random = f
}

func UUID() string {
	b := make([]byte, 16)
	Random.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func DownloadImage(u string) (io.ReadCloser, error) {
	resp, err := http.Get(u)
	return resp.Body, err
}

func Notify(url string, images []image.Image) (err error) {
	data, err := json.Marshal(images)
	if err != nil {
		return
	}
	http.Post(url, "application/json", bytes.NewReader(data))
	return
}
