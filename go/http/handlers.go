package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test_go_app/go/classes"
	"test_go_app/go/db"
	log "test_go_app/go/log"

	"github.com/julienschmidt/httprouter"
)

// Обработка базового запроса
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!")
	log.Print(log.INFO, log.GET_STRING, r.URL, "Hello World!")
}

// Обработка POST запроса получения данных пользователей
func People(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var params db.GetUsersParams
	var filter classes.PeopleFilter
	
	// Чтение тела запроса
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Print(log.ERROR, "Read body error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Парсинг тела запроса
    err = json.Unmarshal(data, &params)
	if err != nil {
		log.Print(log.ERROR, "Unmarshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &filter)
	if err != nil {
		log.Print(log.ERROR, "Unmarshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params.Psql = db.Psql

	// Вызов функции получения пользователей
	db_data, err := db.GetUsers(
		params,
		filter,
	)
	if err != nil {
		return
	}

	data, err = json.MarshalIndent(db_data, "", "    ")
	if err != nil {
		log.Print(log.ERROR, "Marshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	log.Print(log.DEBUG, "%s", data)

	fmt.Fprint(w, string(data))
	log.Print(log.DEBUG, log.POST_STRING, r.URL, fmt.Sprintf("%s %s", filter.String(), params.String()), data)
}