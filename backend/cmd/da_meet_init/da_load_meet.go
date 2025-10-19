package main

import (
	"blreynolds4/event-race-timer/internal/config"
	"blreynolds4/event-race-timer/internal/meets"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	DaID int = iota
	TeamCode
	LastName
	Gender
	Event
	Bib
	Team
	FirstName
	Grade
	Middle
)

func buildAthleteFinder(connectStr string) (meets.AthleteFinder, error) {
	athleteReader, err := meets.NewAthleteReader(connectStr)
	if err != nil {
		return nil, err
	}

	athleteWriter, err := meets.NewAthleteWriter(connectStr)
	if err != nil {
		return nil, err
	}

	return meets.NewAthleteFinder(athleteReader, athleteWriter)
}

func buildMeetFinder(connectStr string) (meets.MeetFinder, error) {
	meetReader, err := meets.NewMeetReader(connectStr)
	if err != nil {
		return nil, err
	}

	meetWriter, err := meets.NewMeetWriter(connectStr)
	if err != nil {
		return nil, err
	}

	return meets.NewMeetFinder(meetReader, meetWriter)
}

func buildRaceFinder(connectStr string) (meets.RaceFinder, error) {
	raceReader, err := meets.NewRaceReader(connectStr)
	if err != nil {
		return nil, err
	}

	raceWriter, err := meets.NewRaceWriter(connectStr)
	if err != nil {
		return nil, err
	}

	return meets.NewRaceFinder(raceReader, raceWriter)
}

func mapRaceName(rawName string) string {
	// Map the DA name to the internal race name format
	cleanName := strings.TrimSpace(rawName)
	if strings.HasPrefix(cleanName, "Combined") {
		return "Combined JV Race"
	}
	return cleanName
}

func main() {
	var claConfigPath string
	var claMeetName string
	var claDaFile string

	flag.StringVar(&claConfigPath, "config", "", "The path to the config file")
	flag.StringVar(&claMeetName, "meetName", "", "The name of the meet")
	flag.StringVar(&claDaFile, "daFile", "", "The da file for event")
	flag.Parse()

	claMeetName = strings.TrimSpace(claMeetName)

	rosterFile, err := os.Create(claMeetName + "_rosters.txt")
	if err != nil {
		fmt.Println("error creating rosters file", err)
		return
	}
	defer rosterFile.Close()

	// Create a default logger with a default log level
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Set the global log level
	}))
	// Set the logger as the default global logger
	slog.SetDefault(logger)

	appConfig := &config.RaceConfig{}
	err = config.LoadConfigData(claConfigPath, appConfig)
	if err != nil {
		fmt.Println("error loading config", err)
		return
	}

	// build meet and athlete finders
	meetFinder, err := buildMeetFinder(appConfig.PgConnect)
	if err != nil {
		fmt.Println("error creating meet finder", err)
		return
	}
	defer meetFinder.Close()

	athleteFinder, err := buildAthleteFinder(appConfig.PgConnect)
	if err != nil {
		fmt.Println("error creating athlete finder", err)
		return
	}
	defer athleteFinder.Close()

	// this will get meet or create it and return it
	meet, err := meetFinder.GetMeet(claMeetName)
	if err != nil {
		fmt.Println("error getting meet", err)
		return
	}

	// create a race finder
	raceFinder, err := buildRaceFinder(appConfig.PgConnect)
	if err != nil {
		fmt.Println("error creating race finder", err)
		return
	}
	defer raceFinder.Close()

	// add competitors and races
	f, err := os.Open(claDaFile)
	if err != nil {
		fmt.Println("error opening event data "+claDaFile, err)
		return
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	// read the file line by line
	// read the header but don't save it
	_, err = csvReader.Read()
	if err != nil {
		return
	}

	// create races and add athletes to them
	lastTeamName := ""
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error reading csv", err)
			return
		}

		// add the event to the meet
		race, err := raceFinder.GetRace(meet, mapRaceName(strings.TrimSpace(record[Event])))
		if err != nil {
			fmt.Println("error getting race", err)
			return
		}

		rawAthlete := meets.Athlete{}
		rawAthlete.DaID = strings.TrimSpace(record[DaID])
		rawAthlete.FirstName = strings.TrimSpace(record[FirstName])
		rawAthlete.LastName = strings.TrimSpace(record[LastName])
		rawAthlete.Gender = strings.TrimSpace(record[Gender])
		rawAthlete.Team = strings.TrimSpace(record[Team])

		year, err := strconv.Atoi(record[Grade])
		if err != nil {
			fmt.Println("err setting grade", err)
			continue
		}
		rawAthlete.Grade = year

		athlete, err := athleteFinder.GetAthlete(rawAthlete)
		if err != nil {
			fmt.Println("error saving athlete", err)
			return
		}

		bib, err := strconv.Atoi(record[Bib])
		if err != nil {
			fmt.Println("err setting bib", err)
			continue
		}

		raceFinder.AddAthlete(race, athlete, bib)
		if lastTeamName != athlete.Team {
			rosterFile.WriteString("\n\n\n\n\n\n\n\n\n\n")
			rosterFile.WriteString(fmt.Sprintf("%s Bib Assignments\n\n", rawAthlete.Team))
		}

		rosterFile.WriteString(fmt.Sprintf("%-20s %-40s %-20s %-5d\n", rawAthlete.Team, rawAthlete.FirstName+" "+rawAthlete.LastName, race.Name, bib))
		lastTeamName = athlete.Team
	}
}
