module github.com/serdarburakguneri/hobby-streamer/backend/lambdas/cmd/trigger_transcode_job

go 1.23.4

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/aws/aws-sdk-go v1.55.7
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../../../pkg/errors

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../../../pkg/logger

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages => ../../../pkg/messages

require github.com/jmespath/go-jmespath v0.4.0 // indirect
