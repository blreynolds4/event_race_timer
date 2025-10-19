package overall

import (
	"blreynolds4/event-race-timer/internal/meets"
	"blreynolds4/event-race-timer/internal/results"
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverallResultsSimple(t *testing.T) {
	athletes := make(meets.AthleteLookup)
	athletes[1] = meets.NewAthlete("JS", "1", "JS", "DAID", 1, "m")
	athletes[10] = meets.NewAthlete("Leb", "1", "Leb", "DAID", 1, "m")
	athletes[11] = meets.NewAthlete("MV", "4", "MV", "DAID", 1, "m")
	athletes[23] = meets.NewAthlete("MV", "10", "MV", "DAID", 1, "m")

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []OverallResult{
		{
			Athlete:    athletes[1],
			Bib:        1,
			Place:      1,
			Finishtime: durationHelper("25m2s"),
		},
		{
			Athlete:    athletes[10],
			Bib:        10,
			Place:      2,
			Finishtime: durationHelper("26m40s"),
		},
		{
			Athlete:    athletes[11],
			Bib:        11,
			Place:      3,
			Finishtime: durationHelper("27m45s"),
		},
		{
			Athlete:    athletes[23],
			Bib:        23,
			Place:      4,
			Finishtime: durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := NewOverallResults(slog.Default())
	err := OVR.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.overallResults)
}

func TestOverallResultsDuplicate(t *testing.T) {
	athletes := make(meets.AthleteLookup)
	athletes[1] = meets.NewAthlete("JS", "1", "JS", "DAID", 1, "m")
	athletes[10] = meets.NewAthlete("Leb", "1", "Leb", "DAID", 1, "m")
	athletes[11] = meets.NewAthlete("MV", "4", "MV", "DAID", 1, "m")
	athletes[23] = meets.NewAthlete("MV", "10", "MV", "DAID", 1, "m")

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	// this way
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))

	// convert these to ^^^
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("22m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []OverallResult{
		{
			Athlete:    athletes[1],
			Bib:        1,
			Place:      1,
			Finishtime: durationHelper("22m2s"),
		},
		{
			Athlete:    athletes[10],
			Bib:        10,
			Place:      2,
			Finishtime: durationHelper("26m40s"),
		},
		{
			Athlete:    athletes[11],
			Bib:        11,
			Place:      3,
			Finishtime: durationHelper("27m45s"),
		},
		{
			Athlete:    athletes[23],
			Bib:        23,
			Place:      4,
			Finishtime: durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := NewOverallResults(slog.Default())
	err := OVR.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.overallResults)
}

func TestOverallResultsError(t *testing.T) {
	athletes := make(meets.AthleteLookup)
	athletes[1] = meets.NewAthlete("JS", "1", "JS", "DAID", 1, "m")
	athletes[10] = meets.NewAthlete("Leb", "1", "Leb", "DAID", 1, "m")
	athletes[11] = meets.NewAthlete("MV", "4", "MV", "DAID", 1, "m")
	athletes[23] = meets.NewAthlete("MV", "10", "MV", "DAID", 1, "m")

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("22m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(meets.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// XCScorer has team results in an array
	OVR := NewOverallResults(slog.Default())
	err := OVR.ScoreResults(context.TODO(), mock)

	assert.Error(t, fmt.Errorf("fail"), err)
}

func durationHelper(d string) time.Duration {
	result, _ := time.ParseDuration(d)
	return result
}

func toMsg(r meets.RaceResult) stream.Message {
	var msg stream.Message
	msgData, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	msg.Data = msgData

	return msg
}
