package classes

import "fmt"

type People struct {
	Passport string   `db:"passport"`
	Name string       `db:"name"`
	Surname string    `db:"surname"`
	Patronymic string `db:"patronymic"`
	Address string    `db:"address"`
}

func (p *People) String() string {
	return fmt.Sprintf("%s %s %s %s %s", p.Passport, p.Name, p.Surname, p.Patronymic, p.Address)
}

type PeopleFilter struct {
	Passport []string   `db:"passport"`
	Name []string       `db:"name"`
	Surname []string    `db:"surname"`
	Patronymic []string `db:"patronymic"`
	Address []string    `db:"address"`
}