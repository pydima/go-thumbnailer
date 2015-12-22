package image

import (
	"io/ioutil"
	"testing"

	"github.com/h2non/bimg"
)

var options bimg.Options = bimg.Options{
	Width:      800,
	Height:     600,
	Enlarge:    true,
	Quality:    95,
	Background: bimg.Color{R: 255, G: 255, B: 255},
	Type:       3, // png
}

func BenchmarkProcessImageGif(b *testing.B) {
	img, err := ioutil.ReadFile("/go/src/github.com/pydima/go-thumbnailer/testdata/gif.gif")
	if err != nil {
		b.Fatalf("Cannot read the test file")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ProcessImage(img, options)
	}
}

func BenchmarkProcessImagePng(b *testing.B) {
	img, err := ioutil.ReadFile("/go/src/github.com/pydima/go-thumbnailer/testdata/png.png")
	if err != nil {
		b.Fatalf("Cannot read the test file")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ProcessImage(img, options)
	}
}

func BenchmarkProcessImageJpg(b *testing.B) {
	img, err := ioutil.ReadFile("/go/src/github.com/pydima/go-thumbnailer/testdata/jpg.jpg")
	if err != nil {
		b.Fatalf("Cannot read the test file")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ProcessImage(img, options)
	}
}

func BenchmarkCreatheThumbnailImagePng(b *testing.B) {
	img, err := ioutil.ReadFile("/go/src/github.com/pydima/go-thumbnailer/testdata/png.png")
	if err != nil {
		b.Fatalf("Cannot read the test file")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		CreateThumbnails(img)
	}
}
