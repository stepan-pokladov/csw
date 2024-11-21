package generator

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stepan-pokladov/csw/server/api"
)

type Generator struct {
	userId uuid.UUID
}

// NewGenerator creates a new generator
func NewGenerator(userId uuid.UUID) *Generator {
	return &Generator{
		userId: userId,
	}
}

// GenerateVisitRecords generates visit records
func (g *Generator) GenerateVisitRecords(startTime time.Time, out chan *api.Visit) int {
	i := 0
	for hour := 0; true; hour += 2 {
		enterTimeMs := startTime.Add(time.Duration(hour) * time.Hour).UnixMilli()
		exitTimeMs := enterTimeMs + (40 * 60 * 1000) // + 40 minutes

		out <- &api.Visit{
			DataVer:       1,
			UserID:        g.userId.String(),
			EnterTime:     enterTimeMs,
			ExitTime:      exitTimeMs,
			AlgorithmType: int64(rand.Intn(7) + 1),
			PoiID:         gofakeit.Int64(),
			Latitude:      gofakeit.Latitude(),
			Longitude:     gofakeit.Longitude(),
		}
		i++
	}
	close(out)
	return i
}

// GenerateActivityRecords generates activity records
func (g *Generator) GenerateActivityRecords(startTime time.Time, out chan *api.Activity) int {
	i := 0
	for hour := 0; true; hour += 2 {
		startTimeMs := startTime.Add(time.Duration(hour) * time.Hour).UnixMilli()
		endTimeMs := startTimeMs + (30 * 60 * 1000) // + 30 minutes
		out <- &api.Activity{
			DataVer:        1,
			UserID:         g.userId.String(),
			ActivityType:   1 << (rand.Intn(6)), // powers of 2 from 1 to 32
			StartTime:      startTimeMs,
			EndTime:        endTimeMs,
			StartLatitude:  gofakeit.Latitude(),
			StartLongitude: gofakeit.Longitude(),
			EndLatitude:    gofakeit.Latitude(),
			EndLongitude:   gofakeit.Longitude(),
		}
		i++
	}
	close(out)
	return i
}
