package meets

import (
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

func TestSaveResult(t *testing.T) {
	// create a meet
	meetWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meet writer: %v", err)
	}
	defer meetWriter.Close()

	raceWriter, err := NewRaceWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create race writer: %v", err)
	}
	defer raceWriter.Close()

	athleteWriter, err := NewAthleteWriter(connectStr)
	assert.Nil(t, err)

	meet := &Meet{
		Name: "Test Meet",
	}
	meet, err = meetWriter.SaveMeet(meet)
	assert.Nil(t, err)
	defer func() {
		meetWriter.DeleteMeet(meet)
	}()

	// create a race
	race := &Race{
		Name: "Test Race",
	}
	race, err = raceWriter.SaveRace(race, meet)
	assert.Nil(t, err)
	assert.Equal(t, race.meet.id, meet.id)

	// create an athlete
	athlete := &Athlete{
		FirstName: "Test",
		LastName:  "Athlete",
		Team:      "Test Team",
		DaID:      "DAID",
		Grade:     1,
		Gender:    "m",
	}
	athlete, err = athleteWriter.SaveAthlete(athlete)
	assert.Nil(t, err)

	defer func() {
		athleteWriter.DeleteAthlete(athlete)
	}()

	// add the athlete to the race
	err = raceWriter.AddAthlete(race, athlete, 1)
	assert.Nil(t, err)

	// save the race result
	raceResult := &RaceResult{
		Bib:         1,
		Athlete:     athlete,
		Place:       1,
		PlaceSource: "manual",
		Time:        time.Second * 10,
	}

	// create race result writer
	resultWriter, err := NewRaceResultWriter(race, connectStr)
	assert.Nil(t, err)
	defer resultWriter.Close()

	savedResult, err := resultWriter.SaveResult(raceResult)
	assert.Nil(t, err)
	assert.Equal(t, raceResult.Bib, savedResult.Bib)
	assert.Equal(t, raceResult.Athlete.id, savedResult.Athlete.id)
	assert.Equal(t, raceResult.Place, savedResult.Place)

	resultReader, err := NewRaceResultReader(race, connectStr)
	assert.Nil(t, err)
	defer resultReader.Close()

	results, err := resultReader.GetRaceResults()
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, raceResult.Bib, results[0].Bib)
	assert.Equal(t, raceResult.Athlete.id, results[0].Athlete.id)
	assert.Equal(t, raceResult.Place, results[0].Place)
	assert.Equal(t, raceResult.Time, results[0].Time)
}
