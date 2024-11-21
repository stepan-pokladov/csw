package helpers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/stepan-pokladov/csw/server/api"
)

// VisitsToInterface converts a slice of visits to a slice of interfaces
func VisitsToInterface(visits []*api.Visit) []interface{} {
	records := make([]interface{}, len(visits))
	for i, v := range visits {
		records[i] = v
	}
	return records
}

// ActivitiesToInterface converts a slice of activities to a slice of interfaces
func ActivitiesToInterface(activities []*api.Activity) []interface{} {
	records := make([]interface{}, len(activities))
	for i, a := range activities {
		records[i] = a
	}
	return records
}

// ToJSONL converts a slice of interfaces to a JSONL byte slice
func ToJSONL(data []interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	for _, record := range data {
		jsonData, err := json.Marshal(record)

		if err != nil {
			return nil, err
		}
		buffer.Write(jsonData)
		buffer.WriteString("\n")
	}
	return buffer.Bytes(), nil
}

// GzipCompress compresses a byte slice using gzip
func GzipCompress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	defer gz.Close()

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	gz.Close()
	return buffer.Bytes(), nil
}
