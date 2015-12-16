package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/h2non/bimg"
)

var colorWhite = bimg.Color{R: 255, G: 255, B: 255}

func createTempConfig(name string) {
	conf := Config{
		ImageParam:  bimg.Options{Width: 800, Height: 600, Background: colorWhite},
		TaskBackend: "Memory",
	}
	j, _ := json.Marshal(conf)
	ioutil.WriteFile(name, j, 0644)
}

func TestReadNonExistentConfig(t *testing.T) {
	c := &Config{}
	err := decodeConfig("/non/existent/path", c)
	if e, ok := err.(*os.PathError); !ok || e.Err != syscall.ENOENT {
		t.Errorf("Unknown exception %s", err)
	}
}

func TestCheckConfigDecoder(t *testing.T) {
	fName := "test_config.json"
	c := &Config{}

	createTempConfig(fName)
	defer os.Remove(fName)

	if err := decodeConfig(fName, c); err != nil {
		t.Error(err)
	}

	if c.ImageParam.Width != 800 ||
		c.ImageParam.Height != 600 ||
		c.ImageParam.Background != colorWhite ||
		c.TaskBackend != "Memory" {
		t.Errorf("Config has been read incorrectly.")
	}
}
