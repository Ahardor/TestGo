package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	log "test_go_app/go/log"

	_ "github.com/lib/pq"

	"test_go_app/go/classes"
)

// Параметры получения пользователей
type GetUsersParams struct {
	Psql *sql.DB
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

	// Формирование запроса при наличии фильтров
	if(len(filter) > 0) {
		query += " WHERE "

		// Перебор полей класса фильтра
		fields := reflect.VisibleFields(reflect.TypeOf(filter[0]))
		values := reflect.ValueOf(filter[0])

		for i, field := range fields {
			if values.Field(i).Len() > 0 {
				if(needAnd) {
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