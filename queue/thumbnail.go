package queue

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

// thumbnail stores information about the image structure including a byte
// array with the JPEG thumbnail.
type thumbnail struct {
	originalWidth  int
	originalHeight int
}

// generateThumbnail opens the image, extracts the dimensions, and generates a
// thumbnail of the photo with the specified constraints.
func generateThumbnail(name, thumbName string, maxWidth, maxHeight int) (*thumbnail, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	originalImage, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	t := resize.Thumbnail(
		uint(maxWidth),
		uint(maxHeight),
		originalImage,
		resize.Lanczos3,
	)
	w, err := os.Create(thumbName)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	if err := jpeg.Encode(w, t, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}
	return &thumbnail{
		originalWidth:  originalImage.Bounds().Dx(),
		originalHeight: originalImage.Bounds().Dy(),
	}, nil
}
