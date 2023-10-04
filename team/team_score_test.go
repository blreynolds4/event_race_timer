package team

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTeamResultsNotEnoughFinishers(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("JS 1", "JS", 1, 1)
	athletes[10] = competitors.NewCompetitor("JS 2", "JS", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)

	// Need to create Results for each athlete
	mock := &results.MockResultSource{
		Results:    make([]results.RaceResult, 0, 4),
		CancelFunc: cancel,
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []teamResult{}

	// XCScorer has team results in an array
	OVR := NewTeamResult()
	err := OVR.ScoreResults(cancelCtx, mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.teamResults)
}

func TestTeamResultsEnoughFinishers(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("MV 1", "MV", 1, 1)
	athletes[10] = competitors.NewCompetitor("MV 2", "MV", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)
	athletes[14] = competitors.NewCompetitor("MV 11", "MV", 1, 1)

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)

	// Need to create Results for each athlete
	mock := &results.MockResultSource{
		Results:    make([]results.RaceResult, 0, 4),
		CancelFunc: cancel,
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[14], Place: 5, Time: durationHelper("37m50s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []teamResult{
		{
			team:            "MV",
			score:           15,
			runnersFinished: 5,
		},
	}

	// XCScorer has team results in an array
	OVR := NewTeamResult()
	err := OVR.ScoreResults(cancelCtx, mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.teamResults)
}

func TestTeamResultsTooManyFinishers(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("MV 1", "MV", 1, 1)
	athletes[10] = competitors.NewCompetitor("MV 2", "MV", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)
	athletes[14] = competitors.NewCompetitor("MV 11", "MV", 1, 1)
	athletes[15] = competitors.NewCompetitor("MV 12", "MV", 1, 1)
	athletes[16] = competitors.NewCompetitor("MV 13", "MV", 1, 1)

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)

	// Need to create Results for each athlete
	mock := &results.MockResultSource{
		Results:    make([]results.RaceResult, 0, 4),
		CancelFunc: cancel,
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[14], Place: 5, Time: durationHelper("37m50s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[15], Place: 6, Time: durationHelper("38m46s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[16], Place: 7, Time: durationHelper("39m50s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []teamResult{
		{
			team:            "MV",
			score:           15,
			runnersFinished: 7,
		},
	}

	// XCScorer has team results in an array
	OVR := NewTeamResult()
	err := OVR.ScoreResults(cancelCtx, mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.teamResults)
}

func durationHelper(d string) time.Duration {
	result, _ := time.ParseDuration(d)
	return result
}
