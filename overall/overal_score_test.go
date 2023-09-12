package overall

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
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
	mock := &results.MockResultSource{
		Results: make([]results.RaceResult, 0, 4),
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []overallResult{
		{
			Athlete: athletes[1],
			Place:   1,
			Ftime:   durationHelper("25m2s"),
		},
		{
			Athlete: athletes[10],
			Place:   2,
			Ftime:   durationHelper("26m40s"),
		},
		{
			Athlete: athletes[11],
			Place:   3,
			Ftime:   durationHelper("27m45s"),
		},
		{
			Athlete: athletes[23],
			Place:   4,
			Ftime:   durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := newOverallResults()
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
	mock := &results.MockResultSource{
		Results: make([]results.RaceResult, 0, 4),
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 1, Athlete: athletes[1], Place: 1, Time: durationHelper("22m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 10, Athlete: athletes[10], Place: 2, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 11, Athlete: athletes[11], Place: 3, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 23, Athlete: athletes[23], Place: 4, Time: durationHelper("37m46s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []overallResult{
		{
			Athlete: athletes[1],
			Place:   1,
			Ftime:   durationHelper("22m2s"),
		},
		{
			Athlete: athletes[10],
			Place:   2,
			Ftime:   durationHelper("26m40s"),
		},
		{
			Athlete: athletes[11],
			Place:   3,
			Ftime:   durationHelper("27m45s"),
		},
		{
			Athlete: athletes[23],
			Place:   4,
			Ftime:   durationHelper("37m46s"),
		},
	}

	// XCScorer has team results in an array
	OVR := newOverallResults()
	err := OVR.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, OVR.overallResults)
}

func durationHelper(d string) time.Duration {
	result, _ := time.ParseDuration(d)
	return result
}
