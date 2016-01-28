package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/h2non/bimg"
)

type Config struct {
	ImageParam      bimg.Options
	TaskBackend     string
	ImageBackend    string
	TmpDir          string
	MediaRoot       string
	Host            string
	Port            int
	ValidExtensions []string
	Thumbnails      map[string][2]int
	Workers         int
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
	err := decodeConfig("/etc/go_thumbnailer/config.json", &Base)
	if err != nil {
		log.Fatalln("Cannot read config. ", err)
	}
}
