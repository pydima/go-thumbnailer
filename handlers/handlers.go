package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

func checkParams(t *tasks.Task) (err error) {
	if t.TaskID == "" {
		t.TaskID = utils.UUID()
	}
	return
}

func CreateThumbnail(w http.ResponseWriter, r *http.Request) {
	var t tasks.Task
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&t); err != nil {
		os.Exit(1)
	}
	tasks.Backend.Put(&t)
	w.WriteHeader(http.StatusCreated)
}
