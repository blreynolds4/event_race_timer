package meets

import "io"

type Athlete struct {
	id        int64
	DaID      string
	FirstName string
	LastName  string
	Team      string
	Grade     int
	Gender    string
}

type RaceAthlete struct {
	Athlete Athlete
	Bib     int
}

type AthleteLookup map[int]*Athlete

type AthleteWriter interface {
	SaveAthlete(athlete *Athlete) (*Athlete, error)
	DeleteAthlete(athlete *Athlete) error
	io.Closer
}

type AthleteReader interface {
	GetAthlete(daID string) (*Athlete, error)
	GetRaceAthlete(r *Race, bib int) (*RaceAthlete, error)
	GetRaceAthletes(r *Race) ([]*RaceAthlete, error)
	io.Closer
}

func NewAthlete(firstName, lastName, team, daID string, grade int, gender string) *Athlete {
	return &Athlete{
		FirstName: firstName,
		LastName:  lastName,
		Team:      team,
		DaID:      daID,
		Grade:     grade,
		Gender:    gender,
	}
}

func (a *Athlete) Name() string {
	return a.FirstName + " " + a.LastName
}
