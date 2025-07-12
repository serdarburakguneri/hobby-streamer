module github.com/serdarburakguneri/hobby-streamer/pkg/auth

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/serdarburakguneri/hobby-streamer/pkg/constants v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/pkg/constants => ../constants
