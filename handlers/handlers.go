package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pydima/go-thumbnailer/tasks"
	"github.com/pydima/go-thumbnailer/utils"
)

func CreateThumbnail(w http.ResponseWriter, r *http.Request) {
	t := tasks.New()
	d := json.NewDecoder(r.Body)

	if err := d.Decode(t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t.TaskID = utils.UUID() // because it might be rewritten by json
	tasks.Backend.Put(t)
	w.WriteHeader(http.StatusCreated)
}
