package meets

import (
	"database/sql"
	"io"
	"log/slog"
	"time"
)

type RaceResultWriter interface {
	SaveResult(rr *RaceResult) (*RaceResult, error)
	io.Closer
}

type RaceResultReader interface {
	GetRaceResults() ([]*RaceResult, error)
	io.Closer
}

func NewRaceResultWriter(r *Race, connectStr string) (RaceResultWriter, error) {
	return buildResultData(r, connectStr)
}

func NewRaceResultReader(r *Race, connectStr string) (RaceResultReader, error) {
	return buildResultData(r, connectStr)
}

type resultData struct {
	db   *sql.DB
	race *Race
}

func (md *resultData) Close() error {
	var err error
	if md.db != nil {
		err = md.db.Close()
		md.db = nil
	}
	return err
}

func buildResultData(r *Race, connectStr string) (*resultData, error) {
	ad := &resultData{
		race: r,
	}

	var err error
	ad.db, err = sql.Open("postgres", connectStr)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	return ad, nil
}

func (rd *resultData) SaveResult(rr *RaceResult) (*RaceResult, error) {
	// race result will be save to athlete_race table
	// the key is race id, bib, athlete id
	slog.Info("Saving race result", "athlete id", rr.Athlete.id, "raceResult", slog.AnyValue(rr))
	query := `
		INSERT INTO athlete_race (race_id, athlete_id, bib, finish_time, place, xc_place, finish_source, place_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (athlete_id, race_id, bib) DO UPDATE SET finish_time = $4, place = $5, xc_place = $6, finish_source = $7, place_source = $8	
	`
	_, err := rd.db.Exec(query, rd.race.id, rr.Athlete.id, rr.Bib, rr.Time.Milliseconds(), rr.Place, rr.XcPlace, rr.FinishSource, rr.PlaceSource)
	if err != nil {
		slog.Error("Error saving race result", slog.String("error", err.Error()))
		return nil, err
	}

	return rr, nil
}

func (rd *resultData) GetRaceResults() ([]*RaceResult, error) {
	query := `
	select 
	  ar.bib,
		a.id,
		a.da_id,
		a.first_name,
		a.last_name,
		a.team,
		a.grade,
		a.gender,
		ar.finish_time, ar.place, ar.xc_place, ar.finish_source, ar.place_source
	FROM athlete a
		JOIN athlete_race ar ON a.id = ar.athlete_id
		inner join race r on ar.race_id = r.id
		inner join meet m on r.meet_id = m.id
	where ar.race_id = $1 and m.id= $2
	ORDER BY ar.place ASC, ar.finish_time ASC
	`
	rows, err := rd.db.Query(query, rd.race.id, rd.race.meet.id)
	if err != nil {
		slog.Error("Error querying results for meet and race", slog.String("meet", rd.race.meet.Name), slog.String("race", rd.race.Name), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	raceResults := make([]*RaceResult, 0, 100)
	for rows.Next() {
		athlete := new(Athlete)
		raceResult := new(RaceResult)
		timeInMillis := int64(0)
		err := rows.Scan(&raceResult.Bib, &athlete.id, &athlete.DaID, &athlete.FirstName, &athlete.LastName, &athlete.Team, &athlete.Grade, &athlete.Gender,
			&timeInMillis, &raceResult.Place, &raceResult.XcPlace, &raceResult.FinishSource, &raceResult.PlaceSource)
		if err != nil {
			slog.Error("Error scanning athlete result row", slog.String("meet", rd.race.meet.Name), slog.String("race", rd.race.Name), slog.String("error", err.Error()))
			continue
		}
		raceResult.Time = time.Duration(timeInMillis) * time.Millisecond

		raceResult.Athlete = athlete
		raceResults = append(raceResults, raceResult)
	}

	slog.Info("Finished getting results for race", slog.String("meet", rd.race.meet.Name), slog.String("race", rd.race.Name), slog.Int("count", len(raceResults)))

	return raceResults, nil
}
