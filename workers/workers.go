package workers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/image/backend"
	"github.com/pydima/go-thumbnailer/models"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

func Run(done <-chan struct{}) {
	tasksChan := make(chan *tasks.Task)
	go func() {
		for {
			t := tasks.Backend.Get()
			tasksChan <- t
		}
	}()

	for {
		select {
		case <-done:
			log.Println("Got signal, stop processing.")
			return
		case t := <-tasksChan:
			fmt.Println("Create task.")
			go process(t)
		}
	}
}

func get_image(is tasks.ImageSource) ([]byte, error) {
	var data []byte

	if is.Path[:4] == "http" {
		return utils.DownloadImage(is.Path)

	} else {
		img, err := os.Open(is.Path)
		if err != nil {
			return data, err
		}
		defer img.Close()
		return ioutil.ReadAll(img)
	}
}

func process(t *tasks.Task) {
	for _, is := range t.Images {
		db_i := models.Image{
			OriginalPath: is.Path,
			Identifier:   is.Identifier,
		}

		if db_i.Exist() {
			log.Println("This image is already exist.")
			return
		}

		s := make(chan []byte, 1)
		go func(is tasks.ImageSource) {
			i, err := get_image(is)
			if err != nil {
				close(s)
				return
			}
			s <- i
		}(is)

		res, ok := <-s
		if !ok {
			continue
		}

		thumbs, err := image.CreateThumbnails(res)
		if err != nil {
			log.Printf("Sorry. %s", err)
			continue
		}

		paths, err := backend.ImageBackend.Save(thumbs)
		if err != nil {
			log.Printf("Shit happens.")
			continue
		}

		db_i.Path = paths[0]

		models.Db.Create(&db_i)
	}
	// var i []image.Image
	// go utils.Notify(t.NotifyUrl, i)

}
