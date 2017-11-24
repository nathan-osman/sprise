package queue

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nathan-osman/sprise/db"
)

// uploadFile attempts to upload a file to a bucket, returning its URL.
func uploadFile(name, filename string, b *db.Bucket) (string, error) {
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
	})
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:    aws.String(b.Name),
		Key:       aws.String(name),
		Body:      f,
		GrantRead: aws.String("http://acs.amazonaws.com/groups/global/AllUsers"),
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
