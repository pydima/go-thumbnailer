package models

import "testing"

func TestPathIfExist(t *testing.T) {
	img := Image{
		OriginalPath: "path",
		Identifier:   "identirier",
	}
	if path := img.PathIfExist(); path != "" {
		t.Errorf("Got: %s, expected: ''", path)
		return
	}

	img.Path = "fsPath"

	if err := Db.Create(&img).Error; err != nil {
		t.Errorf("Cannot create row into the db %s", err.Error())
	}

	img.Path = ""

	if path := img.PathIfExist(); path != "fsPath" {
		t.Errorf("Didn't find the path into db.")
		return
	}

	Db.Where(img).Delete(&Image{})

}
