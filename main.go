package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pydima/go-thumbnailer/config"
	"github.com/pydima/go-thumbnailer/handlers"
	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
	"github.com/pydima/go-thumbnailer/workers"
)

func main() {
	utils.HandleSigTerm()
	http.HandleFunc("/thumbnail", handlers.CreateThumbnail)

	defer tasks.Backend.Close()

	go workers.Run()

	host := config.Base.Host
	port := config.Base.Port
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
}
