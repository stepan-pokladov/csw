package generator

import (
	"github.com/stepan-pokladov/csw/server/api"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVisitRecords(t *testing.T) {
	userID := uuid.New()
	generator := NewGenerator(userID)

	startTime := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	out := make(chan *api.Visit, 100) // Buffered channel to avoid blocking

	// Start generating records
	go generator.GenerateVisitRecords(startTime, out)

	// Collect the first 10 records for testing
	records := make([]*api.Visit, 0)
	for i := 0; i < 10; i++ {
		record := <-out
		records = append(records, record)
	}

	// Assertions
	assert.Equal(t, 10, len(records), "Expected 10 records")
	for i, record := range records {
		assert.Equal(t, int64(1), record.DataVer, "DataVer should always be 1")
		assert.Equal(t, userID.String(), record.UserID, "UserID should match the generator's userId")
		assert.GreaterOrEqual(t, record.Latitude, -90.0, "Latitude should be valid")
		assert.LessOrEqual(t, record.Latitude, 90.0, "Latitude should be valid")
		assert.GreaterOrEqual(t, record.Longitude, -180.0, "Longitude should be valid")
		assert.LessOrEqual(t, record.Longitude, 180.0, "Longitude should be valid")
		assert.Equal(t, record.EnterTime+40*60*1000, record.ExitTime, "ExitTime should be 40 minutes after EnterTime")
		if i > 0 {
			assert.Equal(t, records[i-1].EnterTime+2*60*60*1000, record.EnterTime, "EnterTime should increment by 2 hours")
		}
	}
}

func TestGenerateActivityRecords(t *testing.T) {
	userID := uuid.New()
	generator := NewGenerator(userID)

	startTime := time.Date(2018, time.January, 1, 1, 0, 0, 0, time.UTC)
	out := make(chan *api.Activity, 100) // Buffered channel to avoid blocking

	// Start generating records
	go generator.GenerateActivityRecords(startTime, out)

	// Collect the first 10 records for testing
	records := make([]*api.Activity, 0)
	for i := 0; i < 10; i++ {
		record := <-out
		records = append(records, record)
	}

	// Assertions
	assert.Equal(t, 10, len(records), "Expected 10 records")
	for i, record := range records {
		assert.Equal(t, int64(1), record.DataVer, "DataVer should always be 1")
		assert.Equal(t, userID.String(), record.UserID, "UserID should match the generator's userId")
		assert.Contains(t, []int64{1, 2, 4, 8, 16, 32}, record.ActivityType, "ActivityType should be a power of 2")
		assert.Equal(t, record.StartTime+30*60*1000, record.EndTime, "EndTime should be 30 minutes after StartTime")
		assert.GreaterOrEqual(t, record.StartLatitude, -90.0, "StartLatitude should be valid")
		assert.LessOrEqual(t, record.StartLatitude, 90.0, "StartLatitude should be valid")
		assert.GreaterOrEqual(t, record.StartLongitude, -180.0, "StartLongitude should be valid")
		assert.LessOrEqual(t, record.StartLongitude, 180.0, "StartLongitude should be valid")
		assert.GreaterOrEqual(t, record.EndLatitude, -90.0, "EndLatitude should be valid")
		assert.LessOrEqual(t, record.EndLatitude, 90.0, "EndLatitude should be valid")
		assert.GreaterOrEqual(t, record.EndLongitude, -180.0, "EndLongitude should be valid")
		assert.LessOrEqual(t, record.EndLongitude, 180.0, "EndLongitude should be valid")
		if i > 0 {
			assert.Equal(t, records[i-1].StartTime+2*60*60*1000, record.StartTime, "StartTime should increment by 2 hours")
		}
	}
}
