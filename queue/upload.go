package queue

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
	"github.com/nathan-osman/sprise/db"
)

// uploadFile attempts to upload a file to a bucket.
func uploadFile(name string, r io.Reader, size int64, b *db.Bucket) (string, error) {
	c, err := minio.New(b.Endpoint, b.AccessKey, b.SecretAccessKey, true)
	if err != nil {
		return "", err
	}
	if _, err := c.PutObject(b.Name, name, r, size, minio.PutObjectOptions{}); err != nil {
		return "", err
	}
	return fmt.Sprintf("//%s/%s", b.Endpoint, name), nil
}
