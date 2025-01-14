package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	FileName string `json:"fileName"`
}

type Response struct {
	URL string `json:"url"`
}

func handler(ctx context.Context, req Request) (Response, error) {
	// Get the bucket name from environment variables
	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Println("Environment variable BUCKET_NAME is not set")
		return Response{}, nil
	}

	if req.FileName == "" {
		log.Println("fileName is missing in the request")
		return Response{}, nil
	}

	// Set the AWS region (modify if required)
	region := os.Getenv("BUCKET_REGION")
	if region == "" {
		region = "eu-north-1" // Default region
	}

	// Initialize S3 session with the specified region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return Response{}, err
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Generate a pre-signed URL for PutObject operation
	reqObj, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(req.FileName),
	})
	url, err := reqObj.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate pre-signed URL: %v", err)
		return Response{}, err
	}

	return Response{URL: url}, nil
}

func main() {
	lambda.Start(handler)
}
