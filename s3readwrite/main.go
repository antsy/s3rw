package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, s3Event events.S3Event) (Response, error) {

	event := s3Event.Records[0]

	logEvent(event)

	fmt.Println(fmt.Sprintf("File %s (%d bytes)", event.S3.Object.Key, event.S3.Object.Size))

	sourceFilename := event.S3.Object.Key
	targetFilename := fmt.Sprintf("%s.backup", sourceFilename)
	bucketName := os.Getenv("BUCKET_NAME")

	session := session.Must(session.NewSession())

	serviceClient := s3.New(session, &aws.Config{Region: aws.String("eu-west-1")})
	downloader := s3manager.NewDownloaderWithClient(serviceClient)

	// Buffer to hold file contents in memory
	buffer := aws.NewWriteAtBuffer([]byte{})

	// Read the file from S3
	numBytes, err := downloader.Download(
		buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(sourceFilename),
		})

	if err != nil {
		fmt.Println(fmt.Sprintf("Download failed: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("%d bytes read from S3", numBytes))
	}

	// Convert our byte array to io.reader for writing purposes
	data := bytes.NewReader(buffer.Bytes())

	// Upload input parameters
	uploadParams := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(targetFilename),
		Body:   data,
	}

	uploader := s3manager.NewUploader(session)

	result, err := uploader.Upload(uploadParams)

	if err != nil {
		fmt.Println(fmt.Sprintf("Upload failed: %s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("Upload success: %v", result))
	}

	// Not sure if this makes sense to return from non API gateway Lambda 🤔
	return Response{
		Message: "",
	}, nil
}

func logEvent(event events.S3EventRecord) {
	jsonString, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Fail: ", err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("%s", jsonString))
}

func main() {
	lambda.Start(Handler)
}
