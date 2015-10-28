package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ImageParam struct {
	Width  uint
	Height uint
}

type Config struct {
	ImageParam ImageParam
}

var Base Config

func init() {
	f, _ := os.Open("/etc/fyndiq/config.json")
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&Base); err != nil {
		log.Fatalln("Cannot read config. ", err)
	}
	fmt.Println(Base.ImageParam.Width)
}
