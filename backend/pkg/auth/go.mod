module github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../constants

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../logger

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../errors
