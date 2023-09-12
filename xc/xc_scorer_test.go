package xc

import (
	"blreynolds4/event-race-timer/competitors"
	"blreynolds4/event-race-timer/results"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestXCMeetScore3Teams2Complete(t *testing.T) {
	athletes := make(competitors.CompetitorLookup)
	athletes[1] = competitors.NewCompetitor("JS 1", "JS", 1, 1)
	athletes[2] = competitors.NewCompetitor("JS 2", "JS", 1, 1)
	athletes[3] = competitors.NewCompetitor("JS 3", "JS", 1, 1)
	athletes[4] = competitors.NewCompetitor("JS 4", "JS", 1, 1)
	athletes[5] = competitors.NewCompetitor("JS 5", "JS", 1, 1)
	athletes[6] = competitors.NewCompetitor("JS 6", "JS", 1, 1)
	athletes[7] = competitors.NewCompetitor("JS 7", "JS", 1, 1)
	athletes[8] = competitors.NewCompetitor("JS 8", "JS", 1, 1)
	athletes[9] = competitors.NewCompetitor("JS 9", "JS", 1, 1)
	athletes[10] = competitors.NewCompetitor("Leb 1", "Leb", 1, 1)
	athletes[11] = competitors.NewCompetitor("MV 4", "MV", 1, 1)
	athletes[12] = competitors.NewCompetitor("MV 5", "MV", 1, 1)
	athletes[13] = competitors.NewCompetitor("MV 6", "MV", 1, 1)
	athletes[14] = competitors.NewCompetitor("Leb 2", "Leb", 1, 1)
	athletes[15] = competitors.NewCompetitor("Leb 3", "Leb", 1, 1)
	athletes[16] = competitors.NewCompetitor("Leb 4", "Leb", 1, 1)
	athletes[17] = competitors.NewCompetitor("MV 1", "MV", 1, 1)
	athletes[18] = competitors.NewCompetitor("MV 2", "MV", 1, 1)
	athletes[19] = competitors.NewCompetitor("MV 3", "MV", 1, 1)
	athletes[20] = competitors.NewCompetitor("MV 7", "MV", 1, 1)
	athletes[21] = competitors.NewCompetitor("MV 8", "MV", 1, 1)
	athletes[22] = competitors.NewCompetitor("MV 9", "MV", 1, 1)
	athletes[23] = competitors.NewCompetitor("MV 10", "MV", 1, 1)

	// Need to create Results for each athlete
	mock := results.MockResultSource{
		Results: make([]results.RaceResult, 0, 23),
	}
	mock.Results = append(mock.Results, results.RaceResult{Bib: 17, Athlete: athletes[17], Place: 1, Time: durationHelper("23m29s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[18], Place: 2, Time: durationHelper("24m11s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[1], Place: 3, Time: durationHelper("25m2s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[2], Place: 4, Time: durationHelper("25m10s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[19], Place: 5, Time: durationHelper("25m37s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[10], Place: 6, Time: durationHelper("26m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[2], Place: 7, Time: durationHelper("26m54s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[4], Place: 8, Time: durationHelper("26m57s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[5], Place: 9, Time: durationHelper("27m35s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[11], Place: 10, Time: durationHelper("27m45s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[6], Place: 11, Time: durationHelper("28m44s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[12], Place: 12, Time: durationHelper("28m44s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[13], Place: 13, Time: durationHelper("29m21s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[14], Place: 14, Time: durationHelper("30m40s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[15], Place: 15, Time: durationHelper("31m1s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[7], Place: 16, Time: durationHelper("31m35s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[20], Place: 17, Time: durationHelper("32m27s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[21], Place: 18, Time: durationHelper("34m6s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[22], Place: 19, Time: durationHelper("34m7s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[8], Place: 20, Time: durationHelper("36m22s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[16], Place: 21, Time: durationHelper("36m39s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[9], Place: 22, Time: durationHelper("37m25s")})
	mock.Results = append(mock.Results, results.RaceResult{Bib: 18, Athlete: athletes[23], Place: 23, Time: durationHelper("37m46s")})

	// expected scoring
	// start with just the scores, fill the rest in
	expected := []XCTeamResult{
		{
			Name:      "JS",
			TeamScore: 28,
			Finishers: []XCResult{
				{
					Athlete: athletes[1],
					Score:   3,
				},
				{
					Athlete: athletes[2],
					Score:   4,
				},
				{
					Athlete: athletes[3],
					Score:   6,
				},
				{
					Athlete: athletes[4],
					Score:   7,
				},
				{
					Athlete: athletes[5],
					Score:   8,
				},
				{
					Athlete: athletes[6],
					Score:   10,
				},
				{
					Athlete: athletes[7],
					Score:   13,
				},
				{
					Athlete: athletes[8],
					Score:   0,
				},
				{
					Athlete: athletes[9],
					Score:   0,
				},
			},
		},
		{
			Name:      "MV",
			TeamScore: 28,
			Finishers: []XCResult{
				{
					Athlete: athletes[17],
					Score:   1,
				},
				{
					Athlete: athletes[18],
					Score:   2,
				},
				{
					Athlete: athletes[19],
					Score:   5,
				},
				{
					Athlete: athletes[11],
					Score:   9,
				},
				{
					Athlete: athletes[12],
					Score:   11,
				},
				{
					Athlete: athletes[13],
					Score:   12,
				},
				{
					Athlete: athletes[20],
					Score:   14,
				},
				{
					Athlete: athletes[21],
					Score:   0,
				},
				{
					Athlete: athletes[22],
					Score:   0,
				},
				{
					Athlete: athletes[23],
					Score:   0,
				},
			},
		},
	}

	// XCScorer has team results in an array
	xcScorer := NewXCScorer()
	err := xcScorer.ScoreResults(context.TODO(), mock)
	assert.NoError(t, err)

	assert.Equal(t, expected, xcScorer.Results)
}

func durationHelper(d string) time.Duration {
	result, _ := time.ParseDuration(d)
	return result
}
