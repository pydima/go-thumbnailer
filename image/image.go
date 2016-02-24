package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/h2non/bimg"

	"github.com/pydima/go-thumbnailer/config"
)

var (
	markerJPG = []byte{0xff, 0xd8}
	markerPNG = []byte{0x89, 0x50}
	markerGIF = []byte{0x47, 0x49}
)

type imageType int

// UNKNOWN, JPG, PNG, GIF - constants which define type of image
const (
	UNKNOWN imageType = iota
	JPG
	PNG
	GIF
)

// InvalidImage - error for signaling about problem with image processing
type InvalidImage struct {
	err string
}

func (e InvalidImage) Error() string {
	return e.err
}

func constructName(prefix, ext string) string {
	return fmt.Sprintf("%s.%s", prefix, ext)
}

func checkExtension(n string) error {
	for _, ext := range config.Base.ValidExtensions {
		if strings.HasSuffix(strings.ToLower(n), ext) {
			return nil
		}
	}
	return InvalidImage{fmt.Sprintf("Extension %s is not supported.", filepath.Ext(n))}
}

func imageFormat(img []byte) imageType {
	if len(img) < 2 {
		return UNKNOWN
	}

	switch {
	case bytes.Equal(img[:2], markerJPG):
		return JPG
	case bytes.Equal(img[:2], markerPNG):
		return PNG
	case bytes.Equal(img[:2], markerGIF):
		return GIF
	default:
		return UNKNOWN
	}
}

func imageDimensions(img []byte) (width, height int, err error) {
	r := bytes.NewReader(img)
	conf, _, err := image.DecodeConfig(r)
	return conf.Width, conf.Height, err
}

// vips doesn't support gif natively, so have to convert it with slow standard library
func convertGifToPng(img []byte) ([]byte, error) {
	r := bytes.NewReader(img)
	i, err := gif.Decode(r)
	if err != nil {
		return nil, err
	}

	res := new(bytes.Buffer)
	err = png.Encode(res, i)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func processImage(img []byte, ip config.ImageParam) (res []byte, err error) {
	imgT := imageFormat(img)
	switch imgT {
	case UNKNOWN:
		return nil, fmt.Errorf("got unknown type")
	case GIF:
		img, err = convertGifToPng(img)
		if err != nil {
			return nil, err
		}
	}

	opts := ip.Options
	opts.Type = config.StringToBimgType[ip.Extension]

	meta, err := bimg.Metadata(img)
	if err != nil {
		return nil, err
	}

	if !meta.Alpha {
		opts.Background = bimg.Color{R: 0, G: 0, B: 0}
	}

	return bimg.Resize(img, opts)
}

// CreateThumbnails generates from the source image (function's parameter)
// all thumbnails specified in the config
func CreateThumbnails(original []byte) (map[string][]byte, error) {
	opts := config.Base.Thumbnails[0]
	origThumbnail, err := processImage(original, opts)
	if err != nil {
		return nil, err
	}
	imgs := make(map[string][]byte)
	imgs[constructName(opts.Name, opts.Extension)] = origThumbnail

	for _, v := range config.Base.Thumbnails[1:] {
		img, err := processImage(origThumbnail, v)
		if err != nil {
			return nil, err
		}
		imgs[constructName(v.Name, v.Extension)] = img
	}

	return imgs, nil
}
