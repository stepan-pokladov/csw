package av_processor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/stepan-pokladov/csw/server/api"
	"github.com/stepan-pokladov/csw/server/internal/queue"
)

type Processor struct {
	q queue.QueueProducer
}

func NewAVProcessor(q queue.QueueProducer) *Processor {
	return &Processor{q: q}
}

// Process reads and processes each line of JSONL data
func (p *Processor) Process(jsonl io.Reader, recordType string) error {
	decoder := json.NewDecoder(jsonl)
	lineNumber := 0

	for {
		lineNumber++
		var err error

		switch recordType {
		case "activity":
			var activity api.Activity
			err = decoder.Decode(&activity)
			if err == nil {
				jstr, err := json.Marshal(activity)
				if err != nil {
					log.Printf("Error marshalling activity: %v", err)
				}
				e := p.q.PostMessage("activity", jstr)
				if e != nil {
					return e
				}
			}
		case "visit":
			var visit api.Visit
			err = decoder.Decode(&visit)
			if err == nil {
				jstr, err := json.Marshal(visit)
				if err != nil {
					log.Printf("Error marshalling visit: %v", err)
				}
				e := p.q.PostMessage("visit", jstr)
				if e != nil {
					return e
				}
			}
		default:
			log.Printf("Unknown record type: %s", recordType)
			continue
		}

		if err == io.EOF {
			log.Printf("End of file reached")
			break
		} else if err != nil {
			return fmt.Errorf("error decoding JSONL line %d: %v", lineNumber, err)
		}
	}

	return nil
}
