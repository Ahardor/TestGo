package classes

import "fmt"

type People struct {
	Passport string   `db:"passport" json:"passport"`
	Name string       `db:"name" json:"name"`
	Surname string    `db:"surname" json:"surname"`
	Patronymic string `db:"patronymic" json:"patronymic"`
	Address string    `db:"address" json:"address"`
}

func (p *People) String() string {
	return fmt.Sprintf("%s %s %s %s %s", p.Passport, p.Name, p.Surname, p.Patronymic, p.Address)
}

type PeopleFilter struct {
	Passport []string   `db:"passport" json:"passport"`
	Name []string       `db:"name" json:"name"`
	Surname []string    `db:"surname" json:"surname"`
	Patronymic []string `db:"patronymic" json:"patronymic"`
	Address []string    `db:"address" json:"address"`
}

func (p *PeopleFilter) String() string {
	return fmt.Sprintf("%s %s %s %s %s", p.Passport, p.Name, p.Surname, p.Patronymic, p.Address)
}