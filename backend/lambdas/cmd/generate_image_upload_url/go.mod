module github.com/serdarburakguneri/hobby-streamer/backend/lambdas/cmd/generate_image_upload_url

go 1.21

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/aws/aws-sdk-go v1.50.25
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../../../pkg/logger
