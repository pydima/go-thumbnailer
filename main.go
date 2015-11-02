package main

import (
	"net/http"

	"github.com/pydima/go-thumbnailer/handlers"
	"github.com/pydima/go-thumbnailer/workers"
)

func main() {
	http.HandleFunc("/thumbnail", handlers.CreateThumbnail)

	workers.Run()

	http.ListenAndServe(":8080", nil)
}
