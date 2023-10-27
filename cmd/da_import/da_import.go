package main

import (
	"blreynolds4/event-race-timer/internal/competitors"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	DaID int = iota
	LastName
	Gender
	Event
	Bib
	Team
	FirstName
	Grade
)

func LoadDARunScoreFile(r io.Reader, events map[string]competitors.CompetitorLookup) error {
	csvReader := csv.NewReader(r)

	// read the header but don't save it
	_, err := csvReader.Read()
	if err != nil {
		return err
	}

	combineEvents := make(map[string]string)
	combineEvents["Combined Junior Varsity"] = "Combined JV Race"
	combineEvents["Mixed Junior Varsity"] = "Combined JV Race"

	// read the rest of the records
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		event := record[Event]
		combinedName, found := combineEvents[event]
		if found {
			event = combinedName
		}

		var c competitors.Competitor
		year, err := strconv.Atoi(record[Grade])
		if err != nil {
			fmt.Println("err setting grade", err)
			continue
		}
		c.Grade = year
		c.Team = record[Team]
		c.Name = record[FirstName] + " " + record[LastName]

		bib, err := strconv.Atoi(record[Bib])
		if err != nil {
			fmt.Println("err setting bib", err)
			continue
		}

		// get event from map
		eventAthletes, found := events[event]
		if !found {
			eventAthletes = make(competitors.CompetitorLookup)
			events[event] = eventAthletes
		}

		// add the competitor to the event
		eventAthletes[bib] = &c
	}

	return nil
}

func SaveRostersAndBibs(r io.Reader) error {
	csvReader := csv.NewReader(r)

	// read the header but don't save it
	_, err := csvReader.Read()
	if err != nil {
		return err
	}

	combineEvents := make(map[string]string)
	combineEvents["Combined Junior Varsity"] = "Combined JV Race"
	combineEvents["Mixed Junior Varsity"] = "Combined JV Race"

	// create output
	f, err := os.Create("CAC_bib_roster.txt")
	if err != nil {
		return err
	}

	// read the rest of the records
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		event := record[Event]
		combinedName, found := combineEvents[event]
		if found {
			event = combinedName
		}

		fmt.Fprintf(f, "%-32s %-4s %-36s %-40s\n", record[Team], record[Bib], record[FirstName]+" "+record[LastName], event)
	}

	return nil
}

func main() {
	var claDaFile string

	flag.StringVar(&claDaFile, "daFile", "", "The da file for event")
	flag.Parse()

	f, err := os.Open(claDaFile)
	if err != nil {
		fmt.Println("error opening event data", err)
		return
	}
	defer f.Close()

	events := make(map[string]competitors.CompetitorLookup)

	err = LoadDARunScoreFile(f, events)
	if err != nil {
		fmt.Println("Error reading registration data", err)
	}
	for event, cl := range events {
		eventFile := claDaFile + "-" + strings.ReplaceAll(event, " ", "_") + "-athletes.json"
		cl.Store(eventFile)
	}

	f.Seek(0, 0)
	err = SaveRostersAndBibs(f)
	if err != nil {
		fmt.Println("Error saving roster data", err)
	}

}
