package xc

import (
	"blreynolds4/event-race-timer/internal/results"
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

const resultChunkSize = 25

type XCResult struct {
	Result results.RaceResult
	Score  int16
}

type XCTeamResult struct {
	Name      string
	TeamScore int16
	TotalTime time.Duration
	Top5Avg   time.Duration
	Finishers []*XCResult
	scored    int
}

func NewXCScorer(l *slog.Logger) *XCScorer {
	return &XCScorer{
		logger:  l.With("scorer", "xc"),
		Results: make([]*XCTeamResult, 0),
	}
}

type XCScorer struct {
	logger  *slog.Logger
	Results []*XCTeamResult
}

func (xcs *XCScorer) ScoreResults(ctx context.Context, source results.ResultStream) error {
	teams := make(map[string]*XCTeamResult)

	resultBuffer := make([]results.RaceResult, resultChunkSize)
	resultCount, err := source.GetResults(ctx, resultBuffer)
	if err != nil {
		return err
	}

	// read and get the latest result for each bib
	dedupedResults := make(map[int]results.RaceResult)
	for resultCount > 0 {
		for i := 0; i < resultCount; i++ {
			dedupedResults[resultBuffer[i].Bib] = resultBuffer[i]
			xcs.logger.Debug("xc scorer adding result", "team", resultBuffer[i].Athlete.Team, "bib", resultBuffer[i].Bib, "time", resultBuffer[i].Time, "place", resultBuffer[i].Place)
		}

		resultCount, err = source.GetResults(ctx, resultBuffer)
		if err != nil {
			return err
		}
	}

	// pass one through the results is to build team lists to get finisher counts by team
	// sort the deduped results by place
	sortedResults := make([]*XCResult, 0, len(dedupedResults))
	for _, result := range dedupedResults {
		xcr := new(XCResult)
		xcr.Result = result
		sortedResults = append(sortedResults, xcr)
	}
	sort.SliceStable(sortedResults, func(i, j int) bool { return sortedResults[i].Result.Place < sortedResults[j].Result.Place })

	for i := 0; i < len(sortedResults); i++ {
		var teamResult *XCTeamResult
		teamResult, exists := teams[sortedResults[i].Result.Athlete.Team]
		if !exists {
			teamResult = new(XCTeamResult)
			teamResult.Name = sortedResults[i].Result.Athlete.Team
			teamResult.Finishers = make([]*XCResult, 0)
		}

		teamResult.Finishers = append(teamResult.Finishers, sortedResults[i])
		teams[teamResult.Name] = teamResult
	}

	score := int16(1)
	for i := 0; i < len(sortedResults); i++ {
		// each team has all their finishers set
		// for each finisher assign a score if the team has more 5 or more runners
		// TODO set their score to zero if they are 8th or higher runner
		team := teams[sortedResults[i].Result.Athlete.Team]
		// only give scores to teams of 5+ but less than 8
		if len(team.Finishers) >= 5 && team.scored < 8 {
			team.Finishers[team.scored].Score = score
			if team.scored < 5 {
				team.TeamScore = team.TeamScore + score
			}
			score++
			team.scored++
		}
	}

	// pass two is to sort the teams by score
	xcs.Results = make([]*XCTeamResult, 0)
	sorted := make([]*XCTeamResult, 0, len(teams))
	dnf := make([]*XCTeamResult, 0)
	for _, xcteam := range teams {
		if len(xcteam.Finishers) >= 5 {
			sorted = append(sorted, xcteam)
		} else {
			dnf = append(dnf, xcteam)
		}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		// return i < j, ie i beat j in xc scoring
		if sorted[i].TeamScore == sorted[j].TeamScore {
			// tie between these 2
			// if they have 6 runners, lowest number 6 wins
			if len(sorted[i].Finishers) > 5 &&
				len(sorted[j].Finishers) > 5 {
				if sorted[i].Finishers[5].Result.Place < sorted[j].Finishers[5].Result.Place {
					return true
				} else {
					return false
				}
			} else {
				// if one team doesn't have 6, team with a 6 wins
				return len(sorted[i].Finishers) > len(sorted[j].Finishers)
			}
		}

		// not a tie, low score wins
		return sorted[i].TeamScore < sorted[j].TeamScore
	})

	// sort the dnf teams by finisher count
	sort.SliceStable(dnf, func(i, j int) bool {
		// return i < j, ie i beat j in xc scoring
		// if they have more finishers
		return len(sorted[i].Finishers) > len(sorted[j].Finishers)
	})

	fmt.Printf("%s", "\x1Bc") // clear stdout
	fmt.Printf("\n\n\n")
	fmt.Println("Plc Team                             Score     1    2    3    4    5    6*   7*   8*   9*")
	fmt.Println("=== ================================ =====   ==============================================")
	for i, teamResult := range sorted {
		teamResult.TotalTime = getTeamTime(teamResult)
		teamResult.Top5Avg = getTeamAverage(teamResult)
		xcs.Results = append(xcs.Results, teamResult)
		fmt.Printf("%-3d %-32s %-5d   ", i+1, teamResult.Name, teamResult.TeamScore)
		for f := 0; f < len(teamResult.Finishers); f++ {
			fmt.Printf("%4d ", teamResult.Finishers[f].Score)
		}
		fmt.Printf("\n")
		fmt.Printf("     Total Time: %s\n", time.Unix(0, 0).UTC().Add(teamResult.TotalTime).Format("15:04:05.00"))
		fmt.Printf("        Average: %s\n", time.Unix(0, 0).UTC().Add(teamResult.Top5Avg).Format("04:05.00"))
	}

	for _, dnfTeam := range dnf {
		fmt.Printf("%-3s %-32s\n", "DNP", dnfTeam.Name)
	}

	return nil
}

func getTeamAverage(t *XCTeamResult) time.Duration {
	total := time.Duration(0)
	for i := 0; i < 5; i++ {
		total += t.Finishers[i].Result.Time
	}

	return total / 5
}

func getTeamTime(t *XCTeamResult) time.Duration {
	total := time.Duration(0)
	for i := 0; i < 5; i++ {
		total += t.Finishers[i].Result.Time
	}

	return total
}
