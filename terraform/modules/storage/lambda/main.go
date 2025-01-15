package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Response struct {
	URL string `json:"url"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the body of the request
	var body struct {
		FileName string `json:"fileName"`
	}
	err := json.Unmarshal([]byte(event.Body), &body)
	if err != nil || body.FileName == "" {
		log.Println("fileName is missing or invalid in the request")
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "fileName is required in the request body",
		}, nil
	}

	// Get the bucket name
	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		log.Println("Environment variable BUCKET_NAME is not set")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Server configuration error",
		}, nil
	}

	// Initialize S3 session
	region := os.Getenv("BUCKET_REGION")
	if region == "" {
		region = "eu-north-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Server configuration error",
		}, nil
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Generate a pre-signed URL
	reqObj, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(body.FileName),
	})
	url, err := reqObj.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Failed to generate pre-signed URL: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to generate pre-signed URL",
		}, nil
	}

	// Return the URL as a JSON response
	response := Response{URL: url}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
