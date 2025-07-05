FROM golang:1.22 AS builder

WORKDIR /app
COPY asset-manager/go.mod asset-manager/go.sum ./
RUN go mod download

COPY . .

# Run tests (fail build if tests fail)
RUN cd asset-manager && go test ./... -v

RUN cd asset-manager && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o asset-manager ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/asset-manager/asset-manager .
EXPOSE 8080
ENTRYPOINT ["./asset-manager"]