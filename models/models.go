package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Image struct {
	OriginalPath string `gorm:"primary_key" sql:"type:varchar(255)"`
	Identifier   string `gorm:"primary_key" sql:"type:varchar(255)"`
	Path         string `sql:"type:varchar(255)"`
	CreatedAt    time.Time
}

func (i *Image) Exist() bool {
	return !Db.Where(i).Find(&Image{}).RecordNotFound()
}

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