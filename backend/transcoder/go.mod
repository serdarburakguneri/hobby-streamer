module github.com/serdarburakguneri/hobby-streamer/backend/transcoder

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/service/sqs v1.30.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.47.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.44.0
	github.com/gorilla/mux v1.8.1
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0
)

replace (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ./pkg/constants
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.68 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.20 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)
