package meets

import "io"

type AthleteFinder interface {
	GetAthlete(a Athlete) (*Athlete, error)
	io.Closer
}

type athleteFinderImpl struct {
	reader AthleteReader
	writer AthleteWriter
}

func NewAthleteFinder(rdr AthleteReader, wrtr AthleteWriter) (AthleteFinder, error) {
	return &athleteFinderImpl{
		reader: rdr,
		writer: wrtr,
	}, nil
}

func (a *athleteFinderImpl) GetAthlete(athlete Athlete) (*Athlete, error) {
	foundAthlete, err := a.reader.GetAthlete(athlete.DaID)
	if err != nil {
		return nil, err
	}

	if foundAthlete != nil {
		return foundAthlete, nil
	}

	// Athlete not found, create it
	newAthlete := &Athlete{DaID: athlete.DaID, FirstName: athlete.FirstName, LastName: athlete.LastName, Team: athlete.Team, Grade: athlete.Grade, Gender: athlete.Gender}
	createdAthlete, err := a.writer.SaveAthlete(newAthlete)
	if err != nil {
		return nil, err
	}
	return createdAthlete, nil
}

func (a *athleteFinderImpl) Close() error {
	var err error
	if a.reader != nil {
		err = a.reader.Close()
		a.reader = nil
	}
	if a.writer != nil {
		err = a.writer.Close()
		a.writer = nil
	}
	return err
}
