package overall

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"blreynolds4/event-race-timer/internal/results"
	"blreynolds4/event-race-timer/internal/stream"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverallResultsSimple(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("JS 1", "JS", 1, 1)
	athletes[10] = competitors.NewCompetitor("Leb 1", "Leb", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []OverallResult{
		{
			Athlete:    athletes[1],
			Place:      1,
			Finishtime: durationHelper("25m2s"),
		},
		{
			Athlete:    athletes[10],
			Place:      2,
			Finishtime: durationHelper("26m40s"),
		},
		{
			Athlete:    athletes[11],
			Place:      3,
			Finishtime: durationHelper("27m45s"),
		},
		{
			Athlete:    athletes[23],
			Place:      4,
			Finishtime: durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := NewOverallResults()
	err := OVR.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.overallResults)
}

func TestOverallResultsDuplicate(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("JS 1", "JS", 1, 1)
	athletes[10] = competitors.NewCompetitor("Leb 1", "Leb", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	// this way
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))

	// convert these to ^^^
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("22m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []OverallResult{
		{
			Athlete:    athletes[1],
			Place:      1,
			Finishtime: durationHelper("22m2s"),
		},
		{
			Athlete:    athletes[10],
			Place:      2,
			Finishtime: durationHelper("26m40s"),
		},
		{
			Athlete:    athletes[11],
			Place:      3,
			Finishtime: durationHelper("27m45s"),
		},
		{
			Athlete:    athletes[23],
			Place:      4,
			Finishtime: durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := NewOverallResults()
	err := OVR.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.overallResults)
}

func TestOverallResultsError(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("JS 1", "JS", 1, 1)
	athletes[10] = competitors.NewCompetitor("Leb 1", "Leb", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)

	// Need to create Results for each athlete
	mockEventStream := &stream.MockStream{
		Events: make([]stream.Message, 0),
	}
	mock := results.NewResultStream(mockEventStream)

	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("22m2s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")}))
	mockEventStream.Events = append(mockEventStream.Events, toMsg(results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")}))

	// XCScorer has team results in an array
	OVR := NewOverallResults()
	err := OVR.ScoreResults(context.TODO(), mock)

	assert.Error(t, fmt.Errorf("fail"), err)
}

func durationHelper(d string) time.Duration {
	result, _ := time.ParseDuration(d)
	return result
}

func toMsg(r results.RaceResult) stream.Message {
	var msg stream.Message
	msgData, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	msg.Data = msgData

	return msg
}

func toResult(msg stream.Message) results.RaceResult {
	var rr results.RaceResult
	err := json.Unmarshal(msg.Data, &rr)
	if err != nil {
		panic(err)
	}

	return rr
}
