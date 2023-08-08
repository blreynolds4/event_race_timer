package competitors

type Competitor interface {
	GetName() string
	GetTeam() string
	GetAge() int
	GetGrade() int
}

type CompetitorLookup map[int]Competitor

type competitor struct {
	Name  string
	Team  string
	Age   int
	Grade int
}

func NewCompetitor(name, team string, age, grade int) Competitor {
	return &competitor{
		Name:  name,
		Team:  team,
		Age:   age,
		Grade: grade,
	}
}

func (c *competitor) GetName() string {
	return c.Name
}

func (c *competitor) GetTeam() string {
	return c.Team
}

func (c *competitor) GetAge() int {
	return c.Age
}

func (c *competitor) GetGrade() int {
	return c.Grade
}
