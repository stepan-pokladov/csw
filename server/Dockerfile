# Dockerfile for server
FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o server cmd/server/main.go

CMD ["./server"]
