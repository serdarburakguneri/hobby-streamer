before terraforming:

$go mod init lambda-function
$go mod tidy
$GOOS=linux GOARCH=amd64 go build -o main main.go
$zip main.zip main

after terraforming:

$AWS_PROFILE=sandbox aws lambda invoke \
--function-name generate-presigned-url \
--region eu-north-1 \
--payload '{"fileName": "example-video.mp4"}' \
--cli-binary-format raw-in-base64-out \
response.json


curl -X PUT -T example-video.mp4 "https://....."


