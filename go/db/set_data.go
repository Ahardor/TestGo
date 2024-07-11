package db

import (
	"database/sql"
	"fmt"
	"reflect"

	"test_go_app/go/classes"

	_ "github.com/lib/pq"
)

func DeleteUser(p *sql.DB, passport string) error {
	result, err := p.Exec("DELETE FROM people WHERE passport = $1", passport)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no such user %s", passport)
	}
	return nil
}


func UpdateUser(p *sql.DB, people classes.People) error {
	query := "UPDATE people SET "

	needColon := false

	// Перебор полей класса фильтра
	fields := reflect.VisibleFields(reflect.TypeOf(people))
	values := reflect.ValueOf(people)

	isEmpty := true

	for i, field := range fields {
		if values.Field(i).Len() > 0 {
			// Если нет полей кроме Passport
			if(isEmpty && field.Name != "Passport") {
				isEmpty = false
			}

			if(needColon) {
				query += ", "
			}
			// В тэге db хранится название полей
			query += fmt.Sprintf("%s = '%s'", field.Tag.Get("db"), values.Field(i).Interface())
			needColon = true
		}
	}

	// Проверка на пустоту
	if isEmpty {
		return fmt.Errorf("no fields to update")
	}
	// Формирование запроса
	query += fmt.Sprintf(" WHERE passport = '%s'", people.Passport)

	// Выполнение запроса
	result, err := p.Exec(query)
	if err != nil {
		return err
	}

	// Проверка количества обновленных строк
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}

	// Проверка на наличие обновленных строк
	if rowsAffected == 0 {
		return fmt.Errorf("no such user %s", people.Passport)
	}

	return nil
}

// Добавление пользователя
func InsertUser(p *sql.DB, passport string) error {
	// Выполнение запроса для добавления пользователя
	_, err := p.Exec("INSERT INTO people VALUES($1)", passport)
	if err != nil {
		return err
	}

	return nil
}

// Добавление задачи
func AddTask(p *sql.DB, title string, description string, passport string) error {

	// Выполнение запроса для добавления задачи
	_, err := p.Exec("INSERT INTO tasks (title, description, user_passport) VALUES($1, $2, $3)", title, description, passport)
	if err != nil {
		return err
	}

	return nil
}

// Начало отсчёта времени в задаче
func StartTaskTime(p *sql.DB, taskID string) error {

	// Выполнение запроса для изменения времени
	_, err := p.Exec("UPDATE tasks SET started_at = NOW(), finished_at = NULL WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	return nil
}

// Завершение отсчёта времени в задаче
func FinishTaskTime(p *sql.DB, taskID string) error {

	// Выполнение запроса для изменения времени
	_, err := p.Exec("UPDATE tasks SET finished_at = NOW() WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	return nil
}