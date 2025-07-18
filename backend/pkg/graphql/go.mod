module github.com/serdarburakguneri/hobby-streamer/backend/pkg/graphql

go 1.21

require (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0 // indirect
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../auth

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../constants

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../errors

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../logger
