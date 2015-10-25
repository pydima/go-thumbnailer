package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/utils"
)

type Task struct {
	Path       string
	Delay      bool
	Identifier string
}

func checkParams(t *Task) (err error) {
	if t.Identifier == "" {
		t.Identifier = utils.UUID()
	}
	return
}

func CreateThumbnail(w http.ResponseWriter, r *http.Request) {
	var t Task
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&t); err != nil {
		os.Exit(1)
	}

	if !t.Delay {
		if err := checkParams(&t); err == nil {
			fmt.Fprintf(w, "OK")
		} else {
			log.Fatal("Good buy.")
		}
	}

	var i io.ReadCloser
	if t.Path[:4] == "http" {
		i, _ = utils.DownloadImage(t.Path)
	} else {
		i, _ = utils.ReadImage(t.Path)
	}

	image.ProcessImage(i)
}