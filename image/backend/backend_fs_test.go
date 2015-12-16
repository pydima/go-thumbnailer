package backend

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestCreateTmpDir(t *testing.T) {
	defer cleanUp()
	baseTempDir := filepath.Join(backend.BasePath, "tmp")
	if backend.exists(baseTempDir) {
		t.Errorf("Base temp directory already exists.")
	}

	// create directory when parent directory doesn't exist
	tmpDir, err := backend.createTmpDir()
	if err != nil {
		t.Errorf("Couldn't create directory.")
	}
	os.Remove(tmpDir)

	// create directory when parent directory exists
	tmpDir, err = backend.createTmpDir()
	if err != nil {
		t.Errorf("Couldn't create directory.")
	}
}

// Should always pass image with name 'original.*' in images map.
func TestBaseImage(t *testing.T) {
	// test with 'original.png'
	imgs := map[string][]byte{
		"original.png": []byte("image"),
	}
	_, err := backend.baseImage(imgs)
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	// test with 'original.jpg'
	imgs = map[string][]byte{
		"original.jpg": []byte("image"),
	}
	_, err = backend.baseImage(imgs)
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	// without 'original.*' should get an error
	imgs = map[string][]byte{
		"127x127.png": []byte("image"),
	}
	_, err = backend.baseImage(imgs)
	if err == nil {
		t.Errorf("Should get an error, didn't pass original image.")
	}
}

func TestMoveFiles(t *testing.T) {
	defer cleanUp()

	dst := filepath.Join(backend.BasePath, "new_dir")
	tmpDir, err := backend.createTmpDir()
	if err != nil {
		t.Errorf("Couldn't create directory.")
	}

	// move directory when dst doesn't exist
	err = backend.moveFiles(tmpDir, dst)
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	// move directory when dst exists and empty
	tmpDir, err = backend.createTmpDir()
	if err != nil {
		t.Errorf("Couldn't create directory.")
	}
	err = backend.moveFiles(tmpDir, dst)
	if err != nil {
		t.Errorf("Got error %s", err)
	}

	// move directory when dst exists and not empty
	tmpDir, err = backend.createTmpDir()
	if err != nil {
		t.Errorf("Couldn't create directory.")
	}

	f, err := os.Create(filepath.Join(dst, "file"))
	if err != nil {
		t.Errorf("Couldn't create file %s.", err)
	}
	f.Close()

	err = backend.moveFiles(tmpDir, dst)
	if err == nil {
		t.Errorf("Should get an error if dst is not empty.")
	}
}

func TestImageGC(t *testing.T) {
	defer cleanUp()

	tmpDir, err := backend.createTmpDir()
	if err != nil {
		t.Errorf("Got error %s", err)
	}
	imageGC()
	if !backend.exists(tmpDir) {
		t.Errorf("Directory was deleted.")
	}

	baseTempDir := filepath.Join(backend.BasePath, "tmp")
	tmpDir, err = ioutil.TempDir(baseTempDir, (time.Now().Add(-time.Hour*25)).Format(time.RFC3339)+"_")
	if err != nil {
		t.Errorf("Got error %s", err)
	}
	imageGC()
	if backend.exists(tmpDir) {
		t.Errorf("Directory exists, should have deleted it.")
	}
}
