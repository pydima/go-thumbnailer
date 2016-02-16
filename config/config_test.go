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

func writeTempConfig(name string, conf Config) {
	j, _ := json.Marshal(conf)
	ioutil.WriteFile(name, j, 0644)
}

func createTempConfig() Config {
	ip := make([]ImageParam, 1, 1)
	ip[0] = ImageParam{
		Name:      "original",
		Extension: "png",
	}
	ip[0].Options.Width = 800
	ip[0].Options.Height = 600
	ip[0].Options.Background = colorWhite
	conf := Config{
		Thumbnails:  ip,
		TaskBackend: "Memory",
	}
	return conf
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

	conf := createTempConfig()
	writeTempConfig(fName, conf)
	defer os.Remove(fName)

	if err := decodeConfig(fName, c); err != nil {
		t.Error(err)
	}

	if c.Thumbnails[0].Width != 800 ||
		c.Thumbnails[0].Height != 600 ||
		c.Thumbnails[0].Background != colorWhite ||
		c.TaskBackend != "Memory" {
		t.Errorf("Config has been read incorrectly.")
	}
}

func TestValidateExtensions(t *testing.T) {
	fName := "test_config.json"
	c := &Config{}

	conf := createTempConfig()
	conf.ValidExtensions = []string{"unknown"}
	writeTempConfig(fName, conf)
	defer os.Remove(fName)

	if err := decodeConfig(fName, c); err == nil {
		t.Errorf("Successfully parsed invalid extension.")
	}
}
