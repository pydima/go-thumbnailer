package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/h2non/bimg"
)

type Config struct {
	TaskBackend     string
	ImageBackend    string
	TmpDir          string
	MediaRoot       string
	Host            string
	Port            int
	ValidExtensions []string
	Thumbnails      []ImageParam
	Workers         int
}

type ImageParam struct {
	bimg.Options
	Name      string
	Extension string
}

var SupportedExtensions = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
}

var StringToBimgType = map[string]bimg.ImageType{
	"jpg":  1,
	"jpeg": 1,
	"png":  3,
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
	if err != nil {
		return err
	}
	for _, val := range c.ValidExtensions {
		if _, ok := SupportedExtensions[val]; !ok {
			return fmt.Errorf("unknown extension in valid extension list: %s", val)
		}
	}
	for _, t := range c.Thumbnails {
		if _, ok := StringToBimgType[t.Extension]; !ok {
			return fmt.Errorf("unknown extension in image specification: %s", t.Extension)
		}
	}
	return
}

func init() {
	err := decodeConfig("/etc/go_thumbnailer/config.json", &Base)
	if err != nil {
		log.Fatalln("Cannot read config. ", err)
	}
}
