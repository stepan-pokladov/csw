package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

type KafkaProducer struct {
	kp sarama.SyncProducer
}

// NewKafkaProducer creates a new Kafka queue
func NewKafkaProducer(host string) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	kp, err := sarama.NewSyncProducer([]string{host}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	return &KafkaProducer{kp: kp}
}

// Close closes the Kafka producer
func (q *KafkaProducer) Close() {
	q.kp.Close()
}

// PostMessage sends a message to a Kafka topic
func (q *KafkaProducer) PostMessage(topic string, data []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	}

	partition, offset, err := q.kp.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}
