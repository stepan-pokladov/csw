package queue

type QueueProducer interface {
	PostMessage(topic string, data []byte) error
	Close()
}
