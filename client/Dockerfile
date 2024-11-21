# Dockerfile for client
FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o client cmd/client/main.go

CMD ["./client"]
