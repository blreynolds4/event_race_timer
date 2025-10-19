package competitors

import (
	"encoding/json"
	"os"
)

type CompetitorLookup map[int]*Competitor

type Competitor struct {
	Name  string
	Team  string
	Age   int
	Grade int
}

func NewCompetitor(name, team string, age, grade int) *Competitor {
	return &Competitor{
		Name:  name,
		Team:  team,
		Age:   age,
		Grade: grade,
	}
}

// Implement JSON competitor lookup save and load
func LoadCompetitorLookup(path string, athletes CompetitorLookup) error {
	// read json from the path provided and return the lookup
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(file), &athletes)
}

func (cl CompetitorLookup) Store(path string) error {
	// write json from the path provided and return the lookup
	data, err := json.MarshalIndent(cl, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
