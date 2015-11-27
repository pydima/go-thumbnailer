package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pydima/go-thumbnailer/tasks"
)

func TestCreateUser(t *testing.T) {
	var (
		image_source []tasks.ImageSource
		task         tasks.Task
	)

	image_source = append(image_source, tasks.ImageSource{"http://random_path_to_image.jpg", ""})
	task.Images = image_source
	task.NotifyUrl = "http://localhost:8000/"
	task.TaskID = "Random ID"

	data, err := json.Marshal(task)
	if err != nil {
		os.Exit(1)
	}

	reader := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "", reader)
	w := httptest.NewRecorder()
	http.HandlerFunc(CreateThumbnail).ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected: %d", w.Code)
	}

	task2 := tasks.Backend.Get()
	if task.TaskID != task2.TaskID {
		t.Errorf("Tasks ID's are not the same.")
	}
}
