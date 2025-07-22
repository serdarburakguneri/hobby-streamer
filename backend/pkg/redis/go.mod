module github.com/serdarburakguneri/hobby-streamer/backend/pkg/redis

go 1.21

require (
	github.com/redis/go-redis/v9 v9.3.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
)

replace github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../logger 