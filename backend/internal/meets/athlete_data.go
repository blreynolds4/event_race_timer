package meets

import (
	"database/sql"
	"log/slog"
)

func NewAthleteWriter(connectStr string) (AthleteWriter, error) {
	return buildAthleteData(connectStr)
}

func NewAthleteReader(connectStr string) (AthleteReader, error) {
	return buildAthleteData(connectStr)
}

type athleteData struct {
	db *sql.DB
}

func (md *athleteData) Close() error {
	var err error
	if md.db != nil {
		err = md.db.Close()
		md.db = nil
	}
	return err
}

func buildAthleteData(connectStr string) (*athleteData, error) {
	ad := &athleteData{}

	var err error
	ad.db, err = sql.Open("postgres", connectStr)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	return ad, nil
}

func (ad *athleteData) SaveAthlete(athlete *Athlete) (*Athlete, error) {
	var err error
	if athlete.id == 0 {
		// Insert new athlete
		query := `
		INSERT INTO athlete (da_id, first_name, last_name, team, grade, gender)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
		`
		err = ad.db.QueryRow(query, athlete.DaID, athlete.FirstName, athlete.LastName, athlete.Team, athlete.Grade, athlete.Gender).Scan(&athlete.id)
	} else {
		// Update existing athlete
		query := `
		UPDATE athlete
		SET da_id = $1, first_name = $2, last_name = $3, team = $4, grade = $5, gender = $6
		WHERE id = $7
		`
		_, err = ad.db.Exec(query, athlete.DaID, athlete.FirstName, athlete.LastName, athlete.Team, athlete.Grade, athlete.Gender, athlete.id)
	}
	if err != nil {
		slog.Error("Failed to save athlete", slog.String("error", err.Error()))
		return nil, err
	}
	return athlete, nil
}

func (ad *athleteData) DeleteAthlete(athlete *Athlete) error {
	query := `
		DELETE FROM athlete_race
		WHERE athlete_id = $1
	`
	_, err := ad.db.Exec(query, athlete.id)
	if err != nil {
		slog.Error("Error deleting athlete from all races", slog.String("error", err.Error()))
		return err
	}

	query = `
		DELETE FROM athlete
		WHERE id = $1
	`
	_, err = ad.db.Exec(query, athlete.id)
	if err != nil {
		slog.Error("Error deleting athlete", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (ad *athleteData) GetAthlete(daID string) (*Athlete, error) {

	// Query the database
	row := ad.db.QueryRow("SELECT id, da_id, first_name, last_name, team, grade, gender FROM athlete WHERE da_id = $1", daID)
	athlete := &Athlete{}
	err := row.Scan(&athlete.id, &athlete.DaID, &athlete.FirstName, &athlete.LastName, &athlete.Team, &athlete.Grade, &athlete.Gender)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("No athlete found with da_id", slog.String("da_id", daID))
			return nil, nil
		}
		slog.Error("Error querying athlete by da_id", slog.String("error", err.Error()), slog.String("da_id", daID))
		return nil, err
	}

	return athlete, nil
}

func (ad *athleteData) GetRaceAthletes(r *Race) ([]*RaceAthlete, error) {
	slog.Info("Getting athletes for race", slog.String("race_name", r.Name))

	var raceAthletes []*RaceAthlete
	query := `
	SELECT
		ar.bib,
		a.id,
		a.da_id,
		a.first_name,
		a.last_name,
		a.team,
		a.grade,
		a.gender
	FROM athlete a
		JOIN athlete_race ar ON a.id = ar.athlete_id
		inner join race r on ar.race_id = r.id
		inner join meet m on r.meet_id = m.id
	WHERE m.id = $1 and r.id = $2
	ORDER BY ar.bib`
	rows, err := ad.db.Query(query, r.meet.id, r.id)
	if err != nil {
		slog.Error("Error querying athletes for meet and race", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		athlete := new(RaceAthlete)
		err := rows.Scan(&athlete.Bib, &athlete.Athlete.id, &athlete.Athlete.DaID, &athlete.Athlete.FirstName, &athlete.Athlete.LastName, &athlete.Athlete.Team, &athlete.Athlete.Grade, &athlete.Athlete.Gender)
		if err != nil {
			slog.Error("Error scanning athlete row", slog.String("error", err.Error()))
			return nil, err
		}
		slog.Debug("Adding athlete to race athletes", slog.Int("bib", athlete.Bib), slog.String("name", athlete.Athlete.FirstName+" "+athlete.Athlete.LastName))
		raceAthletes = append(raceAthletes, athlete)
	}

	slog.Info("Finished getting athletes for race", slog.String("race_name", r.Name),
		slog.Int("count", len(raceAthletes)),
		"athletes", raceAthletes)

	return raceAthletes, nil
}

func (ad *athleteData) GetRaceAthlete(r *Race, bib int) (*RaceAthlete, error) {
	slog.Info("Getting athlete for race", slog.String("race_name", r.Name), slog.Int("bib", bib))

	query := `
	SELECT
		ar.bib,
		a.id,
		a.da_id,
		a.first_name,
		a.last_name,
		a.team,
		a.grade,
		a.gender
	FROM athlete a
		JOIN athlete_race ar ON a.id = ar.athlete_id
		inner join race r on ar.race_id = r.id
		inner join meet m on r.meet_id = m.id
	WHERE m.id = $1 and r.id = $2 and ar.bib = $3
	ORDER BY ar.bib`
	row := ad.db.QueryRow(query, r.meet.id, r.id, bib)

	athlete := new(RaceAthlete)
	err := row.Scan(&athlete.Bib, &athlete.Athlete.id, &athlete.Athlete.DaID, &athlete.Athlete.FirstName, &athlete.Athlete.LastName, &athlete.Athlete.Team, &athlete.Athlete.Grade, &athlete.Athlete.Gender)
	if err != nil {
		slog.Error("Error scanning athlete row", slog.String("error", err.Error()))
		return nil, err
	}

	return athlete, nil
}
