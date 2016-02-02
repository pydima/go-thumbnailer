package workers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/image/backend"
	"github.com/pydima/go-thumbnailer/models"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

func run() {
	tasksChan := make(chan *tasks.Task)
	go func() {
		for {
			t := tasks.Backend.Get()
			tasksChan <- t
		}
	}()

	for {
		select {
		case <-utils.STOP:
			log.Println("Got signal, stop processing.")
			return
		case t := <-tasksChan:
			fmt.Println("Create task.")
			process(t)
		}
	}
}

func Run() {
	for x := 0; x < config.Base.Workers; x++ {
		go run()
	}
}

func getImage(is tasks.ImageSource) ([]byte, error) {
	var data []byte

	if len(is.Path) < 4 {
		return data, fmt.Errorf("image path is empty")
	}

	if strings.HasPrefix(is.Path, "http") {
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
	defer tasks.Backend.Complete(t)
	i := make([]string, 0)

	for _, is := range t.Images {
		db_i := models.Image{
			OriginalPath: is.Path,
			Identifier:   is.Identifier,
		}

		if db_i.Exist() {
			log.Println("This image is already exist.")
			return
		}

		res, err := getImage(is)
		if err != nil {
			log.Println(err)
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
		i = append(i, db_i.Path)
	}

	ack := utils.NewAck(t.NotifyUrl, t.ID, i)
	go utils.Notify(ack)

}
