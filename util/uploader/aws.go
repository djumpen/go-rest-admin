package uploader

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/djumpen/go-rest-admin/config"
)

//awsPath - Path to amazon s3 bucket
const awsPath = "https://s3.amazonaws.com/%s/%s"

// NewS3Uploader will return pointer to S3Uploader instance
func NewS3Uploader(fileExt string) *S3Uploader {
	return &S3Uploader{fileExt}
}

// S3Uploader is used to upload file to AWS bucket
type S3Uploader struct {
	FileExt string
}

// Upload will create a new file in AWS bucket from content and return lint to this file
func (u *S3Uploader) Upload(content io.Reader) (url string, err error) {
	aws := config.GetAWS()
	awsCreds := credentials.NewStaticCredentials(aws.ID, aws.Secret, "")

	filename := u.generateFilename()
	if err = u.upload(aws.Region, aws.Bucket, content, awsCreds, filename); err != nil {
		return "", err
	}

	url = fmt.Sprintf(awsPath, aws.Bucket, filename)
	return url, nil
}

// upload is used to upload file to aws
func (u *S3Uploader) upload(region, bucket string, file io.Reader, creds *credentials.Credentials, filename string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	}))

	uploader := s3manager.NewUploader(sess)
	expireTime := time.Now().AddDate(0, 0, +3)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(bucket),
		Key:     aws.String(filename),
		Body:    file,
		Expires: &expireTime,
	})

	return err
}

// generateFilename will return some hashed string with required extension
// this hash is used as filename
func (u *S3Uploader) generateFilename() string {
	hasher := md5.New()

	t, _ := time.Now().MarshalJSON()
	hasher.Write(t)

	return fmt.Sprintf("%s.%s", hex.EncodeToString(hasher.Sum(nil)), u.FileExt)
}
