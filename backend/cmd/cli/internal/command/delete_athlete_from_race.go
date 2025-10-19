package command

import (
	"blreynolds4/event-race-timer/internal/meets"
	"strconv"
)

func NewDeleteAthleteFromRaceCommand(athleteReader meets.AthleteReader, raceReader meets.RaceReader, raceWriter meets.RaceWriter) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			// command line is raceName bib

			// get bib number
			bibNumber, err := strconv.Atoi(args[1])
			if err != nil {
				return false, err
			}

			race, err := raceReader.GetRaceByName(args[0])
			if err != nil || race == nil {
				return false, err
			}

			// need to get the athlete object from race by bib
			athlete, err := athleteReader.GetRaceAthlete(race, bibNumber)
			if err != nil {
				return false, err
			}

			err = raceWriter.RemoveAthlete(race, &athlete.Athlete)
			if err != nil {
				return false, err
			}

			return false, nil
		},
	}
}
