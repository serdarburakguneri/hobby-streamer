module github.com/serdarburakguneri/hobby-streamer/backend/asset-manager

go 1.23.0

toolchain go1.23.10

require (
	github.com/99designs/gqlgen v0.17.74
	github.com/gorilla/mux v1.8.1
	github.com/neo4j/neo4j-go-driver/v5 v5.17.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/vektah/gqlparser/v2 v2.5.30
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../pkg/constants
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
)
