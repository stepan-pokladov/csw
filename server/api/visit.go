package api

import "encoding/json"

func (r *Visit) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Visit struct {
	AlgorithmType int64 `json:"AlgorithmType,omitempty"`
	DataVer       int64 `json:"DataVer,omitempty"`
	// Unix Timestamp in milliseconds
	EnterTime int64 `json:"EnterTime,omitempty"`
	// Unix Timestamp in milliseconds
	ExitTime  int64   `json:"ExitTime,omitempty"`
	Latitude  float64 `json:"Latitude,omitempty"`
	Longitude float64 `json:"Longitude,omitempty"`
	PoiID     int64   `json:"PoiId,omitempty"`
	UserID    string  `json:"UserId,omitempty"`
}
