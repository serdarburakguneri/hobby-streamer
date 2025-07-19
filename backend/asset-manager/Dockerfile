# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

# Copy source code first
COPY . .

# Build the application
RUN cd asset-manager && go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o asset-manager ./cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/asset-manager/asset-manager .
COPY --from=builder /app/asset-manager/config ./config
CMD ["./asset-manager"]