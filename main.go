package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/handlers"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
	"github.com/pydima/go-thumbnailer/workers"
)

var worker = flag.Bool("W", false, "run worker process")

func main() {
	flag.Parse()
	defer tasks.Backend.Close()

	if *worker {
		utils.HandleSigTerm()
		workers.Run()
	} else {
		http.HandleFunc("/thumbnail", handlers.CreateThumbnail)
		host := config.Base.Host
		port := config.Base.Port
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
	}
}
