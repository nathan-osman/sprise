package queue

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nathan-osman/sprise/db"
	"github.com/satori/go.uuid"
)

// uploadFile attempts to upload a file to a bucket, returning its URL.
func uploadFile(filename string, b *db.Bucket) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			b.AccessKey,
			b.SecretAccessKey,
			"",
		),
		Endpoint: aws.String(b.Endpoint),
		Region:   aws.String(b.Region),
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(b.Name),
		Key:         aws.String(uuid.NewV4().String()),
		ContentType: aws.String("image/jpeg"), // TODO: should not be hardcoded
		GrantRead:   aws.String("uri=http://acs.amazonaws.com/groups/global/AllUsers"),
		Body:        f,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
