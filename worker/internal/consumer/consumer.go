package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stepan-pokladov/csw/worker/internal/batch_processor"
)

const (
	groupID      = "g1"
	batchSize    = 100
	batchTimeout = 2 * time.Minute
)

type Consumer struct {
	broker *kafka.Reader
	topic  string
	bp     batch_processor.BatchProcessor
	exitCh chan struct{}
}

// NewConsumer creates a new consumer
func NewConsumer(topic string, bp batch_processor.BatchProcessor, exitCh chan struct{}) *Consumer {
	host := os.Getenv("KAFKA_HOST")
	if host == "" {
		host = "localhost:9092"
	}
	return &Consumer{
		broker: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{host},
			GroupID: groupID,
			Topic:   topic,
		}),
		topic:  topic,
		bp:     bp,
		exitCh: exitCh,
	}
}

// Consume starts consuming messages from the topic
func (c *Consumer) Consume(ctx context.Context) {
	defer c.broker.Close()

	messageChannel := make(chan kafka.Message, 100)
	go func(ctx context.Context) {
		for {
			msg, err := c.broker.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					c.exitCh <- struct{}{}
					return
				}
				log.Printf("Error reading message from topic %s: %v\n", c.topic, err)
				continue
			}
			messageChannel <- msg
		}
	}(ctx)

	var batch []string
	timer := time.NewTimer(batchTimeout)

	fmt.Printf("Started consuming topic: %s\n", c.topic)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Stopping consumer for topic: %s\n", c.topic)
			// Process any remaining batch before exiting
			if len(batch) > 0 {
				err := c.bp.ProcessBatch(c.topic, batch)
				if err != nil {
					log.Printf("Error processing %v batch: %v\n", c.topic, err)
				}
			}
			c.exitCh <- struct{}{}
			return
		case <-timer.C:
			log.Printf("Timeout reached. Processing %v batch...\n", c.topic)
			err := c.bp.ProcessBatch(c.topic, batch)
			if err != nil {
				log.Printf("Error processing %v batch: %v\n", c.topic, err)
			}
			batch = nil
			timer.Reset(batchTimeout)
		case msg := <-messageChannel:
			log.Printf("Message from %s: %s = %s\n", c.topic, string(msg.Key), string(msg.Value))
			batch = append(batch, string(msg.Value))

			// Process batch if it reaches the limit
			if len(batch) >= batchSize {
				err := c.bp.ProcessBatch(c.topic, batch)
				if err != nil {
					log.Printf("Error processing activity batch: %v\n", err)
				}
				batch = nil
				timer.Reset(batchTimeout)
			}
		}
	}
}
