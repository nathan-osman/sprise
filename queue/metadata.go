package queue

import (
	"fmt"
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
		photoDate, _        = x.DateTime()
		photoCameraMake, _  = x.Get(exif.Make)
		photoCameraModel, _ = x.Get(exif.Model)
		photoCamera         string
	)
	if photoCameraMake != nil && photoCameraModel != nil {
		photoCamera = fmt.Sprintf(
			"%s %s",
			photoCameraMake.String(),
			photoCameraModel.String(),
		)
	}
	return &metadata{
		Date:   photoDate,
		Camera: photoCamera,
	}, nil
}
