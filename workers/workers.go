package workers

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pydima/go-thumbnailer/image"
	"github.com/pydima/go-thumbnailer/image/backend"
	"github.com/pydima/go-thumbnailer/models"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

var N = 10

func Run(done <-chan struct{}) {
	for i := 0; i <= N; i++ {
		go process(tasks.Backend, done)
	}
}

func get_image(is tasks.ImageSource) (i io.ReadCloser, err error) {
	if is.Path[:4] == "http" {
		i, err = utils.DownloadImage(is.Path)
	} else {
		i, err = os.Open(is.Path)
	}
	return
}

func process(b tasks.Tasker, done <-chan struct{}) {
	for {
		select {
		case <-done:
			log.Println("Got signal, stop processing.")
		default:
			break
		}

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
			go func(is tasks.ImageSource) {
				i, err := get_image(is)
				if err != nil {
					close(s)
				}
				s <- i
			}(is)

			img, ok := <-s
			if !ok {
				continue
			}

			res, _ := ioutil.ReadAll(img)
			thumbs, err := image.CreateThumbnails(res)
			if err != nil {
				log.Fatal("Sorry.")
			}

			paths, err := backend.ImageBackend.Save(thumbs)
			if err != nil {
				log.Fatal("Shit happens.")
			}

			db_i.Path = paths[0]

			models.Db.Create(&db_i)
		}
		var i []image.Image
		go utils.Notify(t.NotifyUrl, i)
	}
}
