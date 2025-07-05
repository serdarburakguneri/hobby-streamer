FROM golang:1.22 AS builder
WORKDIR /app
COPY transcoder/go.mod transcoder/go.sum ./
RUN go mod download
COPY . .

# Run tests (fail build if tests fail)
RUN cd transcoder && go test ./... -v

RUN cd transcoder && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transcoder-worker ./cmd/worker/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates ffmpeg
WORKDIR /root/
COPY --from=builder /app/transcoder/transcoder-worker .
ENTRYPOINT ["./transcoder-worker"] 