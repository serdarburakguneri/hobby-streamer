after terraforming:

$ffmpeg -f lavfi -i testsrc=duration=5:size=1280x720:rate=30 -c:v libx264 -pix_fmt yuv420p test_video.mp4

$AWS_PROFILE=sandbox aws lambda invoke \
--function-name generate-presigned-url \
--region eu-north-1 \
--payload '{"body": "{\"fileName\": \"test_video.mp4\"}"}' \
--cli-binary-format raw-in-base64-out \
response.json

$cat response.json | jq '.body' | jq -r | jq '.url'

$curl -X PUT -T test_video.mp4 "https://....."




