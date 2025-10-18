package command

import (
	"blreynolds4/event-race-timer/internal/meets"
	"strconv"
)

func NewAddAthleteToRaceCommand(athleteReader meets.AthleteReader, raceReader meets.RaceReader, raceWriter meets.RaceWriter) Command {
	return &noStateCommand{
		CmdFunc: func(args []string) (bool, error) {
			// command line is raceName daid bib

			race, err := raceReader.GetRaceByName(args[0])
			if err != nil || race == nil {
				return false, err
			}

			daid := args[1]

			// get bib number
			bibNumber, err := strconv.Atoi(args[2])
			if err != nil {
				return false, err
			}

			// need to get the athlete object from race by bib
			athlete, err := athleteReader.GetAthlete(daid)
			if err != nil {
				return false, err
			}

			err = raceWriter.AddAthlete(race, athlete, bibNumber)
			if err != nil {
				return false, err
			}

			return false, nil
		},
	}
}
