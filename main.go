package main

import (
	"net/http"

	"github.com/pydima/go-thumbnailer/handlers"
)

func main() {
	http.HandleFunc("/thumbnail", handlers.CreateThumbnail)

	http.ListenAndServe(":8080", nil)
}
