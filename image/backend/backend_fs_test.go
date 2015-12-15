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
	path := "filepath"
	b := backend.exists(path)
	if b {
		t.Errorf("Image already exists.")
	}
}

func TestExist(t *testing.T) {
	defer cleanUp()
	var (
		path []string
		err  error
	)
	image := map[string][]byte{"original.png": []byte("Image")}
	if path, err = backend.Save(image); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	b := backend.exists(path[0])
	if !b {
		t.Errorf("Image doesn't exist.")
	}
}

// Save should check if image already exists
func TestSaveTwice(t *testing.T) {
	defer cleanUp()
	image := map[string][]byte{"original.png": []byte("Image")}
	if _, err := backend.Save(image); err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	if _, err := backend.Save(image); err != nil {
		t.Errorf("Got error %s", err.Error())
	}
}

func TestGenerateDest(t *testing.T) {
	images := make(map[string][]byte)
	images["one"] = []byte("Image one.")
	want := []string{"/root/one"}
	res := backend.generateDest(images, "/root")
	if len(res) != len(want) {
		t.Errorf("Got %s, wanted: %s", res, want)
	}
	if res[0] != want[0] {
		t.Errorf("Got %s, wanted: %s", res, want)
	}
}
