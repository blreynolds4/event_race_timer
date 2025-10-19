package meets

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type meetData struct {
	db *sql.DB
}

func NewMeetReader(connectStr string) (MeetReader, error) {
	return buildMeetData(connectStr)
}

func NewMeetWriter(connectStr string) (MeetWriter, error) {
	return buildMeetData(connectStr)
}

func NewRaceReader(connectStr string) (RaceReader, error) {
	return buildMeetData(connectStr)
}

func NewRaceWriter(connectStr string) (RaceWriter, error) {
	return buildMeetData(connectStr)
}

func (md *meetData) Close() error {
	var err error
	if md.db != nil {
		err = md.db.Close()
		md.db = nil
	}
	return err
}

func buildMeetData(connectStr string) (*meetData, error) {
	md := &meetData{}

	var err error
	md.db, err = sql.Open("postgres", connectStr)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	return md, nil
}

func (md *meetData) GetMeet(name string) (*Meet, error) {
	// Query the database
	row := md.db.QueryRow(`
		SELECT m.id,
			m.name
		FROM meet m
		WHERE m.name = $1`,
		name)

	meet := &Meet{}
	err := row.Scan(&meet.id, &meet.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("No meet found with name", slog.String("name", name))
			return nil, nil
		}
		slog.Error("Error querying meet by name", slog.String("error", err.Error()), slog.String("name", name))
		return nil, err
	}

	return meet, nil
}

func (md *meetData) GetMeets() ([]*Meet, error) {
	// Query the database
	rows, err := md.db.Query(`
		SELECT m.id,
			m.name
		FROM meet m`)
	if err != nil {
		slog.Error("Error querying meets", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	meets := make([]*Meet, 0, 10)
	for rows.Next() {
		meet := new(Meet)
		err := rows.Scan(&meet.id, &meet.Name)
		if err != nil {
			slog.Error("Error scanning meet row", slog.String("error", err.Error()))
			return nil, err
		}
		meets = append(meets, meet)
	}

	return meets, nil
}

func (md *meetData) GetMeetRaces(m *Meet) ([]Race, error) {
	rows, err := md.db.Query(`
		SELECT r.id, r.name
		FROM race r
		WHERE r.meet_id = $1
	`, m.id)
	if err != nil {
		slog.Error("Error querying races", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var races []Race
	for rows.Next() {
		var r Race
		err := rows.Scan(&r.id, &r.Name)
		if err != nil {
			slog.Error("Error scanning row", slog.String("error", err.Error()))
			return nil, err
		}
		// link race to meet
		r.meet = m
		races = append(races, r)
	}

	// save the list of races in the meet
	m.races = races

	return races, nil
}

func (md *meetData) SaveMeet(m *Meet) (*Meet, error) {
	var query string
	if m.id == 0 {
		query = `
        INSERT INTO meet (name)
        VALUES ($1)
				RETURNING id
    `
		err := md.db.QueryRow(query, m.Name).Scan(&m.id)
		if err != nil {
			slog.Error("Error creating meet", slog.String("error", err.Error()))
			return nil, err
		}
	} else {
		query = `
        UPDATE meet
        SET name = $1
        WHERE id = $2
    `
		_, err := md.db.Exec(query, m.Name, m.id)
		if err != nil {
			slog.Error("Error updating meet", slog.String("error", err.Error()))
			return nil, err
		}
	}

	return m, nil
}

func (md *meetData) SaveRace(r *Race, m *Meet) (*Race, error) {
	m, err := md.SaveMeet(m)
	if err != nil {
		return nil, err
	}

	var query string
	if r.id == 0 {
		query = `
        INSERT INTO race (name, meet_id)
        VALUES ($1, $2)
				RETURNING id
    `
		err := md.db.QueryRow(query, r.Name, m.id).Scan(&r.id)
		if err != nil {
			slog.Error("Error creating race", slog.String("error", err.Error()))
			return nil, err
		}
	} else {
		query = `
        UPDATE race
        SET name = $1
        WHERE meet_id = $2 AND id = $3
    `
		_, err := md.db.Exec(query, r.Name, m.id, r.id)
		if err != nil {
			slog.Error("Error updating race", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// make sure meet is linked
	m.AddRace(r)

	return r, nil
}

func (md *meetData) AddAthlete(r *Race, a *Athlete, bib int) error {
	query := `
		INSERT INTO athlete_race (race_id, athlete_id, bib)
		VALUES ($1, $2, $3)
		ON CONFLICT (athlete_id, race_id, bib) DO UPDATE SET bib = $3
	`
	_, err := md.db.Exec(query, r.id, a.id, bib)
	if err != nil {
		slog.Error("Error adding athlete", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (md *meetData) RemoveAthlete(r *Race, a *Athlete) error {
	query := `
		DELETE FROM athlete_race
		WHERE race_id = $1 AND athlete_id = $2
	`
	_, err := md.db.Exec(query, r.id, a.id)
	if err != nil {
		slog.Error("Error removing athlete", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (md *meetData) GetRace(m *Meet, raceName string) (*Race, error) {
	row := md.db.QueryRow("SELECT id, name FROM race WHERE meet_id = $1 AND name = $2", m.id, raceName)
	race := &Race{}
	err := row.Scan(&race.id, &race.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("No race found with name", slog.String("name", raceName))
			return nil, nil
		}
		slog.Error("Error querying race by name", slog.String("error", err.Error()), slog.String("name", raceName))
		return nil, err
	}

	race.meet = m

	return race, nil
}

func (md *meetData) GetRaceByName(raceName string) (*Race, error) {
	rows, err := md.db.Query("SELECT r.id, r.name, m.id, m.name FROM race r join meet m on r.meet_id = m.id WHERE r.name = $1", raceName)
	if err != nil {
		slog.Error("Error querying race by name", slog.String("error", err.Error()), slog.String("name", raceName))
		return nil, err
	}
	defer rows.Close()
	meet := &Meet{}
	race := &Race{}
	foundCount := 0
	for rows.Next() {
		err := rows.Scan(&race.id, &race.Name, &meet.id, &meet.Name)
		if err != nil {
			slog.Error("Error scanning race row", slog.String("error", err.Error()))
			continue
		}
		foundCount++
	}
	if foundCount == 0 {
		slog.Warn("No race found with name", slog.String("name", raceName))
		return nil, nil
	}

	if foundCount > 1 {
		slog.Error("Multiple races found with name", slog.String("name", raceName))
		return nil, fmt.Errorf("multiple races found with name: %s", raceName)
	}

	meet.AddRace(race)

	return race, nil
}

func (md *meetData) DeleteRace(r *Race) error {

	query := `
		DELETE FROM athlete_race
		WHERE race_id = $1
	`
	_, err := md.db.Exec(query, r.id)
	if err != nil {
		slog.Error("Error deleting race athletes", slog.String("error", err.Error()))
		return err
	}

	query = `
		DELETE FROM race
		WHERE id = $1 and meet_id = $2
	`
	_, err = md.db.Exec(query, r.id, r.meet.id)
	if err != nil {
		slog.Error("Error deleting race", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (md *meetData) DeleteMeet(m *Meet) error {
	// delete all the races
	for _, r := range m.races {
		err := md.DeleteRace(&r)
		if err != nil {
			return err
		}
	}

	// delete the meet
	query := `
		DELETE from meet
		WHERE name = $1
	`

	_, err := md.db.Exec(query, m.Name)
	if err != nil {
		slog.Error("Error deleting meet", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func LoadAthleteLookup(connectStr, raceName string, athletes AthleteLookup) error {
	athleteReader, err := NewAthleteReader(connectStr)
	if err != nil {
		return err
	}
	defer athleteReader.Close()

	RaceReader, err := NewRaceReader(connectStr)
	if err != nil {
		return err
	}
	defer RaceReader.Close()

	// get the race by name
	race, err := RaceReader.GetRaceByName(raceName)
	if err != nil {
		return err
	}

	fmt.Println("Loading athletes for race:", raceName)

	// get all athletes for the race
	raceAthletes, err := athleteReader.GetRaceAthletes(race)
	if err != nil {
		return err
	}

	// populate the provided athlete lookup
	for _, ra := range raceAthletes {
		athletes[ra.Bib] = &ra.Athlete
	}

	// success, return nil for error
	return nil
}
