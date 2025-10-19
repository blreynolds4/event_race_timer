package xc

import (
	"blreynolds4/event-race-timer/internal/meets"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

func NewXCTeamScorer(race *meets.Race, l *slog.Logger) *XCTeamScorer {
	return &XCTeamScorer{
		Race:    race,
		logger:  l.With("scorer", "xc"),
		Results: make([]*XCTeamResult, 0),
	}
}

type XCTeamScorer struct {
	Race    *meets.Race
	logger  *slog.Logger
	Results []*XCTeamResult
}

func (xcs *XCTeamScorer) ScoreResults(resultsReader meets.RaceResultReader) error {
	teams := make(map[string]*XCTeamResult)

	// results are returned in place order
	raceResults, err := resultsReader.GetRaceResults()
	if err != nil {
		xcs.logger.Error("ERROR getting race results", "error", err)
		return err
	}

	// pass one through the results is to create an XC result for each result
	xcResults := make([]*XCResult, 0, len(raceResults))
	for _, result := range raceResults {
		xcr := new(XCResult)
		// this should copy the pointer contents into the Result
		xcr.Result = *result
		xcResults = append(xcResults, xcr)
	}

	// group results by team
	for i := 0; i < len(raceResults); i++ {
		var teamResult *XCTeamResult
		teamResult, exists := teams[xcResults[i].Result.Athlete.Team]
		if !exists {
			teamResult = new(XCTeamResult)
			teamResult.Name = xcResults[i].Result.Athlete.Team
			teamResult.Finishers = make([]*XCResult, 0)
		}

		teamResult.Finishers = append(teamResult.Finishers, xcResults[i])
		teams[teamResult.Name] = teamResult
	}

	score := int16(1)
	for i := 0; i < len(xcResults); i++ {
		// each team has all their finishers set
		// for each finisher assign a score if the team has more 5 or more runners
		// TODO set their score to zero if they are 8th or higher runner
		team := teams[xcResults[i].Result.Athlete.Team]
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
		return len(dnf[i].Finishers) > len(dnf[j].Finishers)
	})

	fmt.Printf("%s", "\x1Bc") // clear stdout
	fmt.Printf("Last Updated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
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
		fmt.Printf("%d/5 %-32s\n", len(dnfTeam.Finishers), dnfTeam.Name)
	}

	return nil
}
