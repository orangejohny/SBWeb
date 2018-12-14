package main

import (
	"fmt"
	"log"
	"os"

	"bmstu.codes/developers34/SBWeb/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func initAmazon(m *model.Model) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		fmt.Println("First")
		log.Fatal(err)
	}

	filename := "eee"
	file, _ := os.Open("./AuthReq.PNG")
	uploader := s3manager.NewUploader(sess)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("search-build"),
		Key:    aws.String("/images/" + filename + ".png"),
		Body:   file,
		ACL:    aws.String("public-read"),
	})

	fmt.Println(output.Location)
}

func main() {
	initAmazon(nil)
}
