package config

import (
	"encoding/json"
	"log"
	"os"
)

type ImageParam struct {
	Width  uint
	Height uint
}

type Config struct {
	ImageParam  ImageParam
	TaskBackend string
	Host        string
	Port        int
}

var Base Config

func decodeConfig(path string, c *Config) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(c)
	return
}

func init() {
	err := decodeConfig("/etc/fyndiq/config.json", &Base)
	if err != nil {
		log.Fatalln("Cannot read config. ", err)
	}
}
