package classes

import "fmt"

// Структура данных пользователя в базе
type People struct {
	Passport string   `db:"passport" json:"passport"`
	Name string       `db:"name" json:"name"`
	Surname string    `db:"surname" json:"surname"`
	Patronymic string `db:"patronymic" json:"patronymic"`
	Address string    `db:"address" json:"address"`
}

// Метод для преобразования в строку для People
func (p *People) String() string {
	return fmt.Sprintf("%s %s %s %s %s", p.Passport, p.Name, p.Surname, p.Patronymic, p.Address)
}

// Структура фильтрации пользователей
type PeopleFilter struct {
	Passport []string   `db:"passport" json:"passport"`
	Name []string       `db:"name" json:"name"`
	Surname []string    `db:"surname" json:"surname"`
	Patronymic []string `db:"patronymic" json:"patronymic"`
	Address []string    `db:"address" json:"address"`
}

// Метод для преобразования в строку для PeopleFilter
func (p *PeopleFilter) String() string {
	return fmt.Sprintf("%s %s %s %s %s", p.Passport, p.Name, p.Surname, p.Patronymic, p.Address)
}

// Структура данных задачи
type Task struct {
	ID string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	CreatedAT string `db:"created_at" json:"created_at"`
	StartedAT string `db:"started_at" json:"started_at"`
	FinishedAT string `db:"finished_at" json:"finished_at"`
	UserPassport string `db:"user_passport" json:"user_passport"`
}

// Метод для преобразования в строку для Task
func (t *Task) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s", t.ID, t.Name, t.Description, t.CreatedAT, t.StartedAT, t.FinishedAT, t.UserPassport)
}