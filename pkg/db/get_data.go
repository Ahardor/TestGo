package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	log "test_go_app/pkg/log"

	_ "github.com/lib/pq"

	"test_go_app/pkg/classes"
)

// Параметры получения пользователей
type GetUsersParams struct {
	Psql *sql.DB `swaggerignore:"true"`
	Limit int  `json:"limit"`
	Page int   `json:"page"`
}

func (p *GetUsersParams) String() string {
	return fmt.Sprintf("limit: %d, page: %d", p.Limit, p.Page)
}
// Получение пользователей с фильтрацией
func GetUsers(p GetUsersParams, filter ...classes.PeopleFilter) (result []classes.People, res_err error) {
	if p.Limit == 0 {
		p.Limit = 100
	}

	if p.Page == 0 {
		p.Page = 0
	}

	// Базовый запрос
	query := "SELECT * FROM people"

	needAnd := false
	addedWhere := false

	// Формирование запроса при наличии фильтров
	if(len(filter) > 0) {
		// Перебор полей класса фильтра
		fields := reflect.VisibleFields(reflect.TypeOf(filter[0]))
		values := reflect.ValueOf(filter[0])

		for i, field := range fields {
			if values.Field(i).Len() > 0 {
				if !addedWhere {
					query += " WHERE "
					addedWhere = true
				}
				if needAnd {
					query += " AND "
				}
				// В тэге db хранится название полей
				query += fmt.Sprintf("%s IN ('%s')", field.Tag.Get("db"), strings.Join(values.Field(i).Interface().([]string), "', '"))
				needAnd = true
			}
		}
	}

	// Добавление пагинации
	query += fmt.Sprintf(" LIMIT %d OFFSET %d;", p.Limit, p.Page * p.Limit)

	log.Print(log.DEBUG, "\"%s\"", query)

	// Выполнение запроса и проверка ошибок
	rows, err := p.Psql.Query(query)
	if err != nil {
		log.Print(log.ERROR, "Failed to get users %s", err)
		return nil, err
	}

	// Заполнение результата
	for rows.Next() {
		var p classes.People
		err = rows.Scan(&p.Passport, &p.Name, &p.Surname, &p.Patronymic, &p.Address)
		if err != nil {
			log.Print(log.ERROR, "Failed to get users %s", err)
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

// Получение времени выполнения задач
func GetTime(p *sql.DB, passport string) (time int64, tasksTimes []map[string]int64, err error) {
	tasksTimes = make([]map[string]int64, 0)
	
	// Выполнение запроса
	rows, err := p.Query("SELECT id, started_at, finished_at FROM tasks WHERE user_passport = $1 AND finished_at IS NOT NULL AND started_at IS NOT NULL ORDER BY EXTRACT(EPOCH FROM (finished_at - started_at)) DESC", passport)
	if err != nil {
		return 0, nil, err
	}

	// Заполнение результата
	for rows.Next() {
		var taskId int
		var startedAt sql.NullTime
		var finishedAt sql.NullTime
		// Сканирование полей
		err = rows.Scan(&taskId, &startedAt, &finishedAt)
		if err != nil {
			return 0, nil, err
		}
		// Подсчет времени по отдельным задачам
		tasksTimes = append(tasksTimes, map[string]int64{
			"id": int64(taskId),
			"time": (finishedAt.Time.Unix() - startedAt.Time.Unix()) / 60,
		})
		// Подсчет общего времени
		time = time + (finishedAt.Time.Unix() - startedAt.Time.Unix()) / 60
	}

	return time, tasksTimes, nil
}

// Получение задач пользователя
func GetTasks(p *sql.DB, passport string) (tasks []classes.Task, err error) {
	
	// Выполнение запроса
	rows, err := p.Query("SELECT * FROM tasks WHERE user_passport = $1", passport)
	if err != nil {
		return nil, err
	}

	// Заполнение результата
	for rows.Next() {
		var t classes.Task
		// Сканирование полей
		err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAT, &t.StartedAT, &t.FinishedAT, &t.UserPassport)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}