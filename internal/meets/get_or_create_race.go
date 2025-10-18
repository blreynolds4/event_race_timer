package meets

import "io"

type RaceFinder interface {
	GetRace(m *Meet, raceName string) (*Race, error)
	AddAthlete(r *Race, a *Athlete, bib int) error
	io.Closer
}

type raceFinderImpl struct {
	rdr  RaceReader
	wrtr RaceWriter
}

func NewRaceFinder(rdr RaceReader, wrtr RaceWriter) (RaceFinder, error) {
	return &raceFinderImpl{
		rdr:  rdr,
		wrtr: wrtr,
	}, nil
}

func (rf *raceFinderImpl) GetRace(m *Meet, raceName string) (*Race, error) {
	r, err := rf.rdr.GetRace(m, raceName)
	if err != nil {
		return nil, err
	}
	if r != nil {
		return r, nil
	}

	return rf.wrtr.SaveRace(&Race{Name: raceName, meet: m}, m)
}

func (rf *raceFinderImpl) AddAthlete(r *Race, a *Athlete, bib int) error {
	return rf.wrtr.AddAthlete(r, a, bib)
}

func (rf *raceFinderImpl) Close() error {
	if err := rf.rdr.Close(); err != nil {
		return err
	}
	return rf.wrtr.Close()
}
