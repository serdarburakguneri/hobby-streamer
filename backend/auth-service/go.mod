module github.com/serdarburakguneri/hobby-streamer/backend/auth-service

go 1.23.0

toolchain go1.23.4

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/mux v1.8.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/config v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/security v0.0.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.40.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/config => ../pkg/config
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../pkg/constants
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../pkg/errors
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/security => ../pkg/security
)
