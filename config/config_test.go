package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"syscall"
	"testing"
)

func createTempConfig(name string) {
	conf := Config{
		ImageParam{Width: 800, Height: 600},
		"Memory",
		"localhost",
		8080,
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
		c.TaskBackend != "Memory" {
		t.Errorf("Config has been read incorrectly.")
	}
}
