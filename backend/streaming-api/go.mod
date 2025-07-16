module github.com/serdarburakguneri/hobby-streamer/backend/streaming-api

go 1.23.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/redis/go-redis/v9 v9.11.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0 // indirect
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../pkg/constants

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
