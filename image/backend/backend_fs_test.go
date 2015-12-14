package backend

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var backend FSBackend

func init() {
	path, _ := os.Getwd()
	prefix := "test"
	p, err := ioutil.TempDir(path, prefix)
	if err != nil {
		log.Fatalln("Something went wrong. ", err)
	}
	backend = FSBackend{BasePath: p}
}

func cleanUp() {
	os.RemoveAll(backend.BasePath)
}

func TestNotExist(t *testing.T) {
	hash := "someRandomString"
	b := backend.Exists(hash)
	if b {
		t.Errorf("Image with the hash already exists.")
	}
}

func TestExist(t *testing.T) {
	defer cleanUp()
	hash := "someRandomString"
	var image []byte
	if _, err := backend.Save(image, hash); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	b := backend.Exists(hash)
	if !b {
		t.Errorf("Image with the hash doesn't exist.")
	}
}

// Save should check if image already exists
func TestSaveTwice(t *testing.T) {
	defer cleanUp()
	hash := "someRandomString"
	var image []byte
	if _, err := backend.Save(image, hash); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if _, err := backend.Save(image, hash); err != nil {
		t.Errorf("Got error %s", err.Error())
	}
}
