module github.com/serdarburakguneri/hobby-streamer/backend/asset-manager

go 1.23.0

toolchain go1.23.10

require (
	github.com/99designs/gqlgen v0.17.74
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.0
	github.com/neo4j/neo4j-go-driver/v5 v5.17.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/config v0.0.0-00010101000000-000000000000
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages v0.0.0-00010101000000-000000000000
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/security v0.0.0
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs v0.0.0
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.30
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.36.5 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.26.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.68 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.30.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.20 // indirect
	github.com/aws/smithy-go v1.22.4 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/auth => ../pkg/auth
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/config => ../pkg/config
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants => ../pkg/constants
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors => ../pkg/errors
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/logger => ../pkg/logger
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/messages => ../pkg/messages
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/security => ../pkg/security
	github.com/serdarburakguneri/hobby-streamer/backend/pkg/sqs => ../pkg/sqs
)
