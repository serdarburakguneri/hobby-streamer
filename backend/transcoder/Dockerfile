FROM golang:1.23 AS builder
WORKDIR /app

COPY . .
RUN cd transcoder && go mod download

# Run tests (fail build if tests fail)
RUN cd transcoder && go test ./... -v

RUN cd transcoder && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transcoder-worker ./cmd/worker/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /root/
COPY --from=builder /app/transcoder/transcoder-worker .
COPY --from=builder /app/transcoder/config ./config
ENTRYPOINT ["./transcoder-worker"] 