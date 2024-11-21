package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/stepan-pokladov/csw/worker/internal/batch_processor/file_writer"
	"github.com/stepan-pokladov/csw/worker/internal/consumer"
)

const (
	activityTopic = "activity"
	visitTopic    = "visit"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exitCh := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	activityReader := consumer.NewConsumer(activityTopic, file_writer.NewFileWriter(), exitCh)
	visitReader := consumer.NewConsumer(visitTopic, file_writer.NewFileWriter(), exitCh)

	go activityReader.Consume(ctx)
	go visitReader.Consume(ctx)
	<-exitCh
}
