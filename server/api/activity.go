package api

import "encoding/json"

func (r *Activity) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Activity struct {
	ActivityType int64   `json:"ActivityType,omitempty"`
	DataVer      int64   `json:"DataVer,omitempty"`
	EndLatitude  float64 `json:"EndLatitude,omitempty"`
	EndLongitude float64 `json:"EndLongitude,omitempty"`
	// Unix Timestamp in milliseconds
	EndTime        int64   `json:"EndTime,omitempty"`
	StartLatitude  float64 `json:"StartLatitude,omitempty"`
	StartLongitude float64 `json:"StartLongitude,omitempty"`
	// Unix Timestamp in milliseconds
	StartTime int64  `json:"StartTime,omitempty"`
	UserID    string `json:"UserId,omitempty"`
}
