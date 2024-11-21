package batch_processor

type BatchProcessor interface {
	ProcessBatch(topic string, s []string) error
}
