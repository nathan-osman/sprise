package queue

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// metadata stores information extracted from a photo.
type metadata struct {
	Date   time.Time
	Camera string
}

// parseMetadata reads the EXIF data from an image.
func parseMetadata(name string) (*metadata, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	x, err := exif.Decode(r)
	if err != nil {
		return nil, err
	}
	var (
		date, _      = x.DateTime()
		cameraTag, _ = x.Get(exif.Model)
		camera       string
	)
	if cameraTag != nil {
		camera, _ = cameraTag.StringVal()
	}
	return &metadata{
		Date:   date,
		Camera: camera,
	}, nil
}
