package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Task struct {
	Path       string
	Delay      bool
	Identifier string
}

func CreateThumbnail(w http.ResponseWriter, r *http.Request) {
	var t Task
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&t); err != nil {
		os.Exit(1)
	}

	fmt.Println(t)
}
