package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pydima/go-thumbnailer/tasks"
)

func CreateThumbnail(w http.ResponseWriter, r *http.Request) {
	t := tasks.New()
	d := json.NewDecoder(r.Body)

	if err := d.Decode(t); err != nil {
		os.Exit(1)
	}
	tasks.Backend.Put(t)
	w.WriteHeader(http.StatusCreated)
}
