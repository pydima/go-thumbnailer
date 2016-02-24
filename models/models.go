package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Image contains information with path to the processed image
// and information about the source and time when the image was downloaded
type Image struct {
	OriginalPath string `gorm:"primary_key" sql:"type:varchar(255)"`
	Identifier   string `gorm:"primary_key" sql:"type:varchar(255)"`
	Path         string `sql:"type:varchar(255)"`
	CreatedAt    time.Time
}

// Exist check if image with given URL and identifier exists
func (i *Image) Exist() bool {
	return !Db.Where(i).Find(&Image{}).RecordNotFound()
}

// PathIfExist returns path to the processed thumbnails
// if we have this information into the db
func (i *Image) PathIfExist() string {
	img := new(Image)
	Db.Where(i).Find(img)
	return img.Path
}

// Db is a singleton with connection to the db
var Db gorm.DB

func init() {
	var err error

	Db, err = gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Fatal(err)
	}
	Db.DB()
	Db.DB().Ping()
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
	Db.CreateTable(&Image{})
}
