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

func GetUsers(psql *sql.DB, filter ...classes.PeopleFilter) (result []classes.People, res_err error) {
	query := "SELECT * FROM people"

	needAnd := false

	if(len(filter) > 0) {
		query += " WHERE "
		fields := reflect.VisibleFields(reflect.TypeOf(filter[0]))
		values := reflect.ValueOf(filter[0])

		for i, field := range fields {
			if values.Field(i).Len() > 0 {
				if(needAnd) {
					query += " AND "
				}
				query += fmt.Sprintf(" %s IN (%s) ", field.Tag.Get("db"), strings.Join(values.Field(i).Interface().([]string), ", "))
				needAnd = true
			}
		}

		log.Print(log.DEBUG, query)
		// if len(filter[0].Passport) > 0 {
		// 	query += fmt.Sprintf("passport IN ('%s') ", strings.Join(filter[0].Passport, ", "))
		// 	needAnd = true
		// }

		// if len(filter[0].Name) > 0 {
		// 	if needAnd {
		// 		query += "AND "
		// 	}
		// 	query += fmt.Sprintf("name IN ('%s') ", strings.Join(filter[0].Name, ", "))
		// 	needAnd = true
		// }

		// if len(filter[0].Surname) > 0 {
		// 	if needAnd {
		// 		query += "AND "
		// 	}
		// 	query += fmt.Sprintf("surname IN ('%s') ", strings.Join(filter[0].Surname, ", "))
		// 	needAnd = true
		// }

		// if len(filter[0].Patronymic) > 0 {
		// 	if needAnd {
		// 		query += "AND "
		// 	}
		// 	query += fmt.Sprintf("patronymic IN ('%s') ", strings.Join(filter[0].Patronymic, ", "))
		// 	needAnd = true
		// }
	}

	rows, err := psql.Query("SELECT * FROM people")
	if err != nil {
		log.Print(log.ERROR, "Failed to get users %s", err)
		return nil, err
	}

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