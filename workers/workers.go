package workers

import (
	"io"
	"log"
	"os"

	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/models"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

var N = 10

func Run() {
	for i := 0; i <= N; i++ {
		go process(tasks.Backend)
	}
}

func get_image(is tasks.ImageSource) (i io.ReadCloser) {
	if is.Path[:4] == "http" {
		i, _ = utils.DownloadImage(is.Path)
	} else {
		i, _ = os.Open(is.Path)
	}
	return
}

func process(b tasks.Tasker) {
	t := b.Get()
	for _, is := range t.Images {

		db_i := models.Image{
			OriginalPath: is.Path,
			Identifier:   is.Identifier,
		}

		if db_i.Exist() {
			log.Println("This image is already exist.")
			break
		}

		s := make(chan io.ReadCloser, 1)
		go func() {
			i := get_image(is)
			s <- i
		}()

		path, err := image.ProcessImage(<-s)
		if err != nil {
			log.Fatal("Sorry.")
		}

		db_i.Path = path

		models.Db.Create(&db_i)
	}
	var i []image.Image
	go utils.Notify(t.NotifyUrl, i)
}
