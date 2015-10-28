package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/models"
	"github.com/pydima/go-thumbnailer/utils"
)

type ImageSource struct {
	Path       string
	Identifier string
}

type Task struct {
	Images []ImageSource
	Delay  bool
	TaskID string
}

func checkParams(t *Task) (err error) {
	if t.TaskID == "" {
		t.TaskID = utils.UUID()
	}
	return
}

func get_image(is ImageSource) (i io.ReadCloser) {
	if is.Path[:4] == "http" {
		i, _ = utils.DownloadImage(is.Path)
	} else {
		i, _ = utils.ReadImage(is.Path)
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

	for _, is := range t.Images {

		db_i := models.Image{
			OriginalPath: is.Path,
			Identifier:   is.Identifier,
		}

		if db_i.Exist() {
			fmt.Println("This image is already exist.")
			return
		}

		s := make(chan io.ReadCloser, 1)
		go func(s chan<- io.ReadCloser) {
			i := get_image(is)
			s <- i
		}(s)

		path, err := image.ProcessImage(<-s)
		if err != nil {
			log.Fatal("Sorry.")
		}

		db_i.Path = path

		models.Db.Create(&db_i)
	}

}
