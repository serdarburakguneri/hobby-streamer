FROM golang:1.23 AS builder
WORKDIR /app

COPY . .
RUN cd streaming-api && go mod download

# Build the application
RUN cd streaming-api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o streaming-api ./cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/streaming-api/streaming-api .
COPY --from=builder /app/streaming-api/config ./config
CMD ["./streaming-api"]
