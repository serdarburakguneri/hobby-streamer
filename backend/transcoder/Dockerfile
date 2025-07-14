FROM golang:1.23 AS builder
WORKDIR /app
COPY . .

# Ensure go.sum is up to date
RUN cd transcoder && go mod tidy

# Run tests (fail build if tests fail)
RUN cd transcoder && go test ./... -v

RUN cd transcoder && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transcoder-worker ./cmd/worker/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /root/
COPY --from=builder /app/transcoder/transcoder-worker .
ENTRYPOINT ["./transcoder-worker"] 