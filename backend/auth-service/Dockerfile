FROM golang:1.23-alpine as builder
WORKDIR /app

COPY . .
RUN cd auth-service && go mod download

# Run tests (fail build if tests fail)
RUN cd auth-service && go test ./... -v

RUN cd auth-service && go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth-service/main .
COPY --from=builder /app/auth-service/config ./config
EXPOSE 8080
CMD ["./main"] 