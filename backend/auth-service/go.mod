module github.com/serdarburakguneri/hobby-streamer/backend/auth-service

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0
)

replace (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ./pkg/auth
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ./pkg/constants
)
