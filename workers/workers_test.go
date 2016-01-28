package workers

import (
	"testing"

	"github.com/pydima/go-thumbnailer/tasks"
)

func TestGetImage(t *testing.T) {
	is := tasks.ImageSource{Path: ""}
	_, err := getImage(is)
	if err == nil {
		t.Errorf("Should have returned an error.")
	}
}
