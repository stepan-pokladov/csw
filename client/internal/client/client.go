package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stepan-pokladov/csw/client/pkg/generator"
	"github.com/stepan-pokladov/csw/client/pkg/generator/helpers"
	"github.com/stepan-pokladov/csw/server/api"
)

type Client struct {
	visitUrl       string
	activityUrl    string
	reportInterval time.Duration
	generator      *generator.Generator
}

// NewClient creates a new client
func NewClient(visitUrl, activityUrl string, reportInterval time.Duration) *Client {
	return &Client{
		visitUrl:       visitUrl,
		activityUrl:    activityUrl,
		reportInterval: reportInterval,
		generator:      generator.NewGenerator(uuid.New()),
	}
}

// Run starts the client
func (c *Client) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go c.startPostingVisits()
	wg.Add(1)
	c.startPostingActivities()
	wg.Wait()
}

// startPostingVisits generates and posts visit records
func (c *Client) startPostingVisits() {
	vc := make(chan *api.Visit)
	visitStartTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	go c.generator.GenerateVisitRecords(visitStartTime, vc)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var batch []*api.Visit
		var currentDay int
		for v := range vc {
			// On day change post the batch
			if currentDay != time.UnixMilli(v.EnterTime).UTC().YearDay() {
				err := c.postBatch("visits", helpers.VisitsToInterface(batch))
				if err != nil {
					log.Printf("Error posting visits batch: %v\n", err)
				}
				currentDay = time.UnixMilli(v.EnterTime).UTC().YearDay()
				if len(batch) == 0 {
					continue
				}
				batch = nil
				time.Sleep(c.reportInterval)
			}
			batch = append(batch, v)
		}
	}()

	wg.Wait()
}

// startPostingActivities generates and posts activity records
func (c *Client) startPostingActivities() {
	ac := make(chan *api.Activity)
	activityStartTime := time.Date(2018, 1, 1, 1, 0, 0, 0, time.UTC)
	go c.generator.GenerateActivityRecords(activityStartTime, ac)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var batch []*api.Activity
		var currentDay int
		for a := range ac {
			// On day change post the batch
			if currentDay != time.UnixMilli(a.StartTime).UTC().YearDay() {
				err := c.postBatch("activities", helpers.ActivitiesToInterface(batch))
				if err != nil {
					log.Printf("Error posting visits batch: %v\n", err)
				}
				currentDay = time.UnixMilli(a.StartTime).UTC().YearDay()
				if len(batch) == 0 {
					continue
				}
				batch = nil
				time.Sleep(c.reportInterval)
			}
			batch = append(batch, a)
		}
	}()

	wg.Wait()
}

// postBatch creates a JSONL file from the records and compresses it before sending it to the server
func (c *Client) postBatch(batchType string, records []interface{}) error {
	if len(records) == 0 {
		return nil
	}

	jsonl, err := helpers.ToJSONL(records)
	if err != nil {
		return err
	}
	gzipData, err := helpers.GzipCompress(jsonl)
	if err != nil {
		return err
	}

	var url string
	switch batchType {
	case "visits":
		url = c.visitUrl
	case "activities":
		url = c.activityUrl
	default:
		return fmt.Errorf("unknown batch type: %s", batchType)
	}

	return send(url, gzipData)
}

// send sends the compressed JSONL data to the server
func send(endpoint string, data []byte) error {
	log.Printf("Sending report to %s\n", endpoint)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/jsonl")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send report: %s", body)
	}

	fmt.Printf("Report sent to %s: %s\n", endpoint, body)
	return nil
}
