package meets

import "io"

type Meet struct {
	id    int64
	Name  string
	races []Race
}

type Race struct {
	id   int64
	Name string
	meet *Meet
}

type MeetReader interface {
	GetMeet(name string) (*Meet, error)
	GetMeets() ([]*Meet, error)
	GetMeetRaces(m *Meet) ([]Race, error)
	io.Closer
}

type MeetWriter interface {
	SaveMeet(m *Meet) (*Meet, error)
	DeleteMeet(m *Meet) error
	io.Closer
}

type RaceReader interface {
	GetRace(m *Meet, raceName string) (*Race, error)
	GetRaceByName(raceName string) (*Race, error)
	io.Closer
}

type RaceWriter interface {
	SaveRace(r *Race, m *Meet) (*Race, error)
	AddAthlete(r *Race, a *Athlete, bib int) error
	RemoveAthlete(r *Race, a *Athlete) error
	DeleteRace(r *Race) error
	io.Closer
}

func (m *Meet) AddRace(race *Race) *Race {
	// check if race already exists
	for _, r := range m.races {
		if r.Name == race.Name {
			r.meet = m
			return &r
		}
	}

	// not found in meet, add it
	race.meet = m
	m.races = append(m.races, *race)
	return race
}
