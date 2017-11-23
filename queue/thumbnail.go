package queue

import (
	"bytes"
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
	data           []byte
}

// generateThumbnail opens the image, extracts the dimensions, and generates a
// thumbnail of the photo with the specified constraints.
func generateThumbnail(name string, maxWidth, maxHeight int) (*thumbnail, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	originalImage, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	t := resize.Thumbnail(
		uint(maxWidth),
		uint(maxHeight),
		originalImage,
		resize.Lanczos3,
	)
	b := &bytes.Buffer{}
	if err := jpeg.Encode(b, t, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}
	return &thumbnail{
		originalWidth:  originalImage.Bounds().Dx(),
		originalHeight: originalImage.Bounds().Dy(),
		data:           b.Bytes(),
	}, nil
}
