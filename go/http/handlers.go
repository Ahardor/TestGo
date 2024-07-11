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
func GetPeople(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// Получение переменной ссылки подключения к базе данных
	params.Psql = db.Psql

	// Вызов функции получения пользователей
	db_data, err := db.GetUsers(
		params,
		filter,
	)
	if err != nil {
		return
	}

	// Преобразование полученных данных в JSON
	data, err = json.MarshalIndent(db_data, "", "    ")
	if err != nil {
		log.Print(log.ERROR, "Marshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправка ответа
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	log.Print(log.DEBUG, "%s", data)

	fmt.Fprint(w, string(data))

	log.Print(log.DEBUG, log.POST_STRING, r.URL, fmt.Sprintf("%s %s", filter.String(), params.String()), data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var params struct {
		Passport string `json:"passport"`
	}

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

	// Вызов функции удаления пользователя
	err = db.DeleteUser(db.Psql, params.Passport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to delete user %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.DELETE_STRING, r.URL, params.Passport)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Переменная для хранения параметров запроса
	var params classes.People

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

	if params.Passport == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Passport is empty")
		return
	}

	// Вызов функции обновления пользователя
	err = db.UpdateUser(db.Psql, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to update user %s", err.Error())
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Переменная для хранения параметров запроса
	var params struct {
		Passport string `json:"passport"`
	}

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

	// Вызов функции создания пользователя
	err = db.InsertUser(db.Psql, params.Passport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to create user %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.POST_STRING, r.URL, params.Passport, "OK")
}