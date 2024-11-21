package main

import (
	"github.com/stepan-pokladov/csw/server/internal/client"
	"github.com/stepan-pokladov/csw/server/internal/queue/kafka"
	"os"
)

func main() {
	host := os.Getenv("KAFKA_HOST")
	if host == "" {
		host = "localhost:9092"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	q := kafka.NewKafkaProducer(host)
	service := client.NewService(q)
	service.Run(port)
}
