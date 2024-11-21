package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stepan-pokladov/csw/server/internal/report_processor"
)

// HandlerForRoute determines which type of record to process based on the route
func HandlerForRoute(recordType string, rp report_processor.ReportProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/jsonl" {
			http.Error(w, "Invalid Content-Type. Expected 'text/jsonl'", http.StatusBadRequest)
			return
		}

		log.Printf("Starting JSONL processing...\n")
		if err := rp.Process(r.Body, recordType); err != nil {
			log.Printf("Error processing JSONL data: %v\n", err)
			http.Error(w, fmt.Sprintf("Error processing JSONL data: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File processed successfully"))
	}
}
