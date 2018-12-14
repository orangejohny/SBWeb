package s3

import (
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// Config is config for AWS S3 image manager
type Config struct {
	Bucket string `json:"Bucket"`
	ACL    string `json:"ACL"`
	Region string `json:"Region"`
}

// S3 is a struct that implements model.IM interface
type S3 struct {
	sess   *session.Session
	bucket string
	acl    string
	region string
}

// InitS3 initiates connection to S3
func InitS3(cfg Config) (*S3, error) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil, errors.New("No required environment variables setted")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, err
	}

	r := &S3{
		sess:   sess,
		bucket: cfg.Bucket,
		acl:    cfg.ACL,
		region: cfg.Region,
	}

	return r, nil
}

// IsExist checks if such key exists in the bucket
func (s *S3) IsExist(key string) bool {
	svc := s3.New(s.sess)
	_, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Range:  aws.String("bytes=0-1"),
	})

	if err != nil {
		return false
	}
	return true
}

// UploadImage uploads body as key
func (s *S3) UploadImage(key string, body io.Reader) (string, error) {
	uploader := s3manager.NewUploader(s.sess)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
		ACL:    aws.String(s.acl),
	})
	if err != nil {
		return "", err
	}
	return output.Location, nil
}

// DeleteImage deletes such key from the bucket
func (s *S3) DeleteImage(key string) error {
	svc := s3.New(s.sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(key)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}

// DownloadImage downloads image from AWS
func (s *S3) DownloadImage(key string) error {
	downloader := s3manager.NewDownloader(s.sess)
	file, _ := os.Create("." + key)
	defer file.Close()
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		})
	return err
}
