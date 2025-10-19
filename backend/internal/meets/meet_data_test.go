package meets

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

var connectStr = "postgres://testdb:testdb@localhost:5433/testdb?sslmode=disable"

func TestCloseAClosedReader(t *testing.T) {
	// Create a meetData instance
	md, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetData: %v", err)
	}
	md.Close()
	err = md.Close()
	assert.Nil(t, err)
}

func TestGetMeetFromDb(t *testing.T) {
	// Create a meetData instance (nothing cached)
	md, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetData: %v", err)
	}
	defer md.Close()

	// Insert a test meet into the database
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec("INSERT INTO meet (id, name) VALUES ($1, $2)", 1, "Test Meet")
	if err != nil {
		t.Fatalf("Failed to insert test meet: %v", err)
	}

	// defer cleanup
	defer func() {
		_, err := db.Exec("DELETE FROM meet WHERE id = $1", 1)
		if err != nil {
			t.Fatalf("Failed to delete test meet: %v", err)
		}
	}()

	// Test: Retrieve the meet by name
	meet, err := md.GetMeet("Test Meet")
	assert.Nil(t, err)
	assert.NotNil(t, meet)
	assert.Equal(t, int64(1), meet.id)
	assert.Equal(t, "Test Meet", meet.Name)
	assert.Equal(t, 0, len(meet.races))
}

func TestGetMeetFromDbNotFound(t *testing.T) {
	// Create a meetData instance
	md, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetData: %v", err)
	}
	defer md.Close()

	// Test: Retrieve the meet by name
	meet, err := md.GetMeet("Test Meet")
	assert.Nil(t, err)
	assert.Nil(t, meet)
}

func TestGetMeetListFromDb(t *testing.T) {
	// Create a meetData instance (nothing cached)
	md, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetData: %v", err)
	}
	defer md.Close()

	// Insert a test meet into the database
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec("INSERT INTO meet (id, name) VALUES ($1, $2)", 1, "Test Meet 1")
	if err != nil {
		t.Fatalf("Failed to insert test meet: %v", err)
	}
	_, err = db.Exec("INSERT INTO meet (id, name) VALUES ($1, $2)", 2, "Test Meet 2")
	if err != nil {
		t.Fatalf("Failed to insert test meet 2: %v", err)
	}

	// defer cleanup
	defer func() {
		_, err := db.Exec("DELETE FROM meet WHERE id = $1", 1)
		if err != nil {
			t.Fatalf("Failed to delete test meet: %v", err)
		}
		_, err = db.Exec("DELETE FROM meet WHERE id = $1", 2)
		if err != nil {
			t.Fatalf("Failed to delete test meet: %v", err)
		}
	}()

	// Test: Retrieve the meet by name
	meets, err := md.GetMeets()
	assert.Nil(t, err)
	assert.NotNil(t, meets)
	assert.Equal(t, 2, len(meets))
	assert.Equal(t, int64(1), meets[0].id)
	assert.Equal(t, "Test Meet 1", meets[0].Name)
	assert.Equal(t, 0, len(meets[0].races))
	assert.Equal(t, int64(2), meets[1].id)
	assert.Equal(t, "Test Meet 2", meets[1].Name)
	assert.Equal(t, 0, len(meets[1].races))
}

func TestCreateNewDeleteMeet(t *testing.T) {
	// Insert a test meet into the database
	mWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetWriter: %v", err)
	}
	defer mWriter.Close()

	// save a new meet
	meet := &Meet{Name: "Test Meet"}
	saved, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, saved)
	assert.NotZero(t, saved.id)
	assert.Equal(t, "Test Meet", saved.Name)

	// Test: Delete the meet
	err = mWriter.DeleteMeet(meet)
	assert.Nil(t, err)

	mReader, err := NewMeetReader(connectStr)
	assert.Nil(t, err)
	defer mReader.Close()

	// Test: Retrieve the meet by name
	meet, err = mReader.GetMeet(meet.Name)
	assert.Nil(t, err)
	assert.Nil(t, meet)
}

func TestUpdateDeleteMeet(t *testing.T) {
	// Insert a test meet into the database
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec("INSERT INTO meet (id, name) VALUES ($1, $2)", 1, "Test Meet")
	if err != nil {
		t.Fatalf("Failed to insert test meet: %v", err)
	}

	// defer cleanup
	defer func() {
		_, err := db.Exec("DELETE FROM meet WHERE id = $1", 1)
		if err != nil {
			t.Fatalf("Failed to delete test meet: %v", err)
		}
	}()

	mReader, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetReader: %v", err)
	}
	defer mReader.Close()

	meet, err := mReader.GetMeet("Test Meet")
	assert.Nil(t, err)
	assert.NotNil(t, meet)
	assert.Equal(t, int64(1), meet.id)
	assert.Equal(t, "Test Meet", meet.Name)

	mWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetWriter: %v", err)
	}
	defer mWriter.Close()

	// save an updated meet
	meet.Name = "Test Update"
	saved, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, "Test Update", saved.Name)
	assert.Equal(t, int64(1), saved.id)

	// Test: Delete the meet
	err = mWriter.DeleteMeet(meet)
	assert.Nil(t, err)

	// Test: Retrieve the meet by name
	meet, err = mReader.GetMeet(meet.Name)
	assert.Nil(t, err)
	assert.Nil(t, meet)
}

func TestCreateUpdateAndDeleteRace(t *testing.T) {
	mWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetWriter: %v", err)
	}
	defer mWriter.Close()

	// save a new meet
	meet := &Meet{Name: "Test Meet"}
	saved, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, saved)
	assert.NotZero(t, saved.id)
	assert.Equal(t, "Test Meet", saved.Name)
	defer func() {
		err := mWriter.DeleteMeet(saved)
		assert.Nil(t, err)
	}()

	raceWriter, err := NewRaceWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create raceWriter: %v", err)
	}
	defer raceWriter.Close()

	// test: create a race in the meet
	race := &Race{Name: "Test Race", meet: saved}
	savedRace, err := raceWriter.SaveRace(race, saved)
	assert.Nil(t, err)
	assert.NotNil(t, savedRace)
	assert.NotZero(t, savedRace.id)
	assert.Equal(t, "Test Race", savedRace.Name)
	assert.Equal(t, saved, savedRace.meet)

	raceReader, err := NewRaceReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create raceReader: %v", err)
	}
	defer raceReader.Close()

	// test: retrieve the race by name
	retrievedRace, err := raceReader.GetRace(saved, savedRace.Name)
	assert.Nil(t, err)
	assert.NotNil(t, retrievedRace)
	assert.Equal(t, savedRace.id, retrievedRace.id)
	assert.Equal(t, savedRace.Name, retrievedRace.Name)
	assert.Equal(t, saved, retrievedRace.meet)

	retrievedRace.Name = "Test Race Updated"
	updatedRace, err := raceWriter.SaveRace(retrievedRace, saved)
	assert.Nil(t, err)
	assert.Equal(t, retrievedRace.id, updatedRace.id)
	assert.Equal(t, "Test Race Updated", updatedRace.Name)

	err = raceWriter.DeleteRace(savedRace)
	assert.Nil(t, err)

	retrievedRace, err = raceReader.GetRace(saved, savedRace.Name)
	assert.Nil(t, err)
	assert.Nil(t, retrievedRace)

	deletedByMeetRace := &Race{Name: "Test Race Deleted by meet", meet: saved}
	lastRace, err := raceWriter.SaveRace(deletedByMeetRace, saved)
	assert.Nil(t, err)
	saved.AddRace(deletedByMeetRace)

	// Test: Delete the meet
	err = mWriter.DeleteMeet(meet)
	assert.Nil(t, err)

	retrievedRace, err = raceReader.GetRace(saved, lastRace.Name)
	assert.Nil(t, err)
	assert.Nil(t, retrievedRace)

}

func TestGetMeetRaces(t *testing.T) {
	// save a meet
	mWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetWriter: %v", err)
	}
	defer mWriter.Close()

	raceWriter, err := NewRaceWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create raceWriter: %v", err)
	}
	defer raceWriter.Close()

	meet := &Meet{Name: "Test Meet"}
	saved, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, saved)
	assert.NotZero(t, saved.id)
	assert.Equal(t, "Test Meet", saved.Name)

	// add 2 races to the meet
	race1 := &Race{Name: "Test Race 1", meet: saved}
	race2 := &Race{Name: "Test Race 2", meet: saved}
	_, err = raceWriter.SaveRace(race1, saved)
	assert.Nil(t, err)
	_, err = raceWriter.SaveRace(race2, saved)
	assert.Nil(t, err)

	// retrieve all races for the meet
	meetReader, err := NewMeetReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetReader: %v", err)
	}
	defer meetReader.Close()
	races, err := meetReader.GetMeetRaces(saved)
	assert.Nil(t, err)
	assert.NotNil(t, races)
	assert.Len(t, races, 2)
	assert.Equal(t, "Test Race 1", races[0].Name)
	assert.Equal(t, "Test Race 2", races[1].Name)

	// link races to meet so they get deleted
	for _, r := range races {
		saved.AddRace(&r)
	}

	mWriter.DeleteMeet(saved)
}

func TestGetRaceByName(t *testing.T) {
	// save a meet
	mWriter, err := NewMeetWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create meetWriter: %v", err)
	}
	defer mWriter.Close()

	raceWriter, err := NewRaceWriter(connectStr)
	if err != nil {
		t.Fatalf("Failed to create raceWriter: %v", err)
	}
	defer raceWriter.Close()

	meet := &Meet{Name: "Test Meet"}
	saved, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, saved)
	assert.NotZero(t, saved.id)
	assert.Equal(t, "Test Meet", saved.Name)

	// add race to the meet
	race1 := &Race{Name: "Test Race 1", meet: saved}
	_, err = raceWriter.SaveRace(race1, saved)
	assert.Nil(t, err)

	// retrieve all races for the meet
	raceReader, err := NewRaceReader(connectStr)
	if err != nil {
		t.Fatalf("Failed to create raceReader: %v", err)
	}
	defer raceReader.Close()
	race, err := raceReader.GetRaceByName(race1.Name)
	assert.Nil(t, err)
	assert.NotNil(t, race)
	assert.Equal(t, "Test Race 1", race.Name)
	assert.NotNil(t, race.meet)
	assert.Equal(t, race.meet.id, saved.id)

	mWriter.DeleteMeet(saved)
}

func TestAddAthleteToRace(t *testing.T) {
	mWriter, err := NewMeetWriter(connectStr)
	assert.Nil(t, err)
	assert.NotNil(t, mWriter)
	defer mWriter.Close()

	raceWriter, err := NewRaceWriter(connectStr)
	assert.Nil(t, err)
	assert.NotNil(t, raceWriter)
	defer raceWriter.Close()

	// create meet
	meet := &Meet{Name: "Test Meet"}
	savedMeet, err := mWriter.SaveMeet(meet)
	assert.Nil(t, err)
	assert.NotNil(t, savedMeet)
	assert.NotZero(t, savedMeet.id)
	assert.Equal(t, "Test Meet", savedMeet.Name)

	// create a race
	race := &Race{Name: "Test Race", meet: savedMeet}
	savedRace, err := raceWriter.SaveRace(race, savedMeet)
	assert.Nil(t, err)
	assert.NotNil(t, savedRace)
	assert.NotZero(t, savedRace.id)
	assert.Equal(t, "Test Race", savedRace.Name)
	assert.Equal(t, savedMeet, savedRace.meet)
	assert.Equal(t, 1, len(savedMeet.races))

	athleteWriter, err := NewAthleteWriter(connectStr)
	assert.Nil(t, err)
	assert.NotNil(t, athleteWriter)
	defer athleteWriter.Close()

	// create an athlete
	athlete := &Athlete{DaID: "xxx", FirstName: "Test", LastName: "Athlete", Team: "Test", Grade: 12, Gender: "M"}
	savedAthlete, err := athleteWriter.SaveAthlete(athlete)
	assert.Nil(t, err)
	assert.NotNil(t, savedAthlete)
	assert.NotZero(t, savedAthlete.id)
	assert.Equal(t, athlete.DaID, savedAthlete.DaID)
	assert.Equal(t, athlete.FirstName, savedAthlete.FirstName)
	assert.Equal(t, athlete.LastName, savedAthlete.LastName)
	assert.Equal(t, athlete.Team, savedAthlete.Team)
	assert.Equal(t, athlete.Grade, savedAthlete.Grade)
	assert.Equal(t, athlete.Gender, savedAthlete.Gender)

	// add the athlete to the race
	bib := 123
	err = raceWriter.AddAthlete(savedRace, savedAthlete, bib)
	assert.Nil(t, err)

	athleteReader, err := NewAthleteReader(connectStr)
	assert.Nil(t, err)
	assert.NotNil(t, athleteReader)
	defer athleteReader.Close()

	// read the athletes for the race
	raceAthletes, err := athleteReader.GetRaceAthletes(savedRace)
	assert.Nil(t, err)
	assert.NotNil(t, raceAthletes)
	assert.Len(t, raceAthletes, 1)
	assert.Equal(t, savedAthlete.DaID, raceAthletes[0].Athlete.DaID)
	assert.Equal(t, bib, raceAthletes[0].Bib)

	// delete the meet
	err = mWriter.DeleteMeet(savedMeet)
	assert.Nil(t, err)

	// delete the athlete
	err = athleteWriter.DeleteAthlete(savedAthlete)
	assert.Nil(t, err)
}
