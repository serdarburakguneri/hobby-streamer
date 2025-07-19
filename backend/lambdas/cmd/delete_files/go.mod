module github.com/serdarburakguneri/hobby-streamer/backend/lambdas/cmd/delete_files

go 1.23.0

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.53.8
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../../../pkg/errors

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../../../pkg/logger

require github.com/jmespath/go-jmespath v0.4.0 // indirect
