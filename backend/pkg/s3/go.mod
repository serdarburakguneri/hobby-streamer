module github.com/serdarburakguneri/hobby-streamer/backend/pkg/s3

go 1.23.0

toolchain go1.23.4

require (
	github.com/aws/aws-sdk-go v1.53.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../logger
replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../errors

require github.com/jmespath/go-jmespath v0.4.0 // indirect
