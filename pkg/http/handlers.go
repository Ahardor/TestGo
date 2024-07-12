package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test_go_app/pkg/classes"
	"test_go_app/pkg/db"
	log "test_go_app/pkg/log"

	_ "test_go_app/docs"

	"github.com/julienschmidt/httprouter"
)

// Обработка базового запроса
// @Summary Index
// @Tags index
// @Description Базовые слова для отображения
// @ID index
// @Success 200 {string} string "Hello World!"
// @Failure 400 {string} string "Bad Request"
// @Router / [get]
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!")
	log.Print(log.INFO, log.GET_STRING, r.URL, "Hello World!")
}

// Обработка POST запроса получения данных пользователей
// @Summary GetPeople
// @Tags people
// @Description Данные пользователей
// @ID getPeople
// @Accept  json
// @Produce  json
// @Param params body db.GetUsersParams true "Параметры запроса"
// @Param filters body classes.PeopleFilter true "Параметры запроса"
// @Success 200 {object} []classes.People
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /getPeople [post]
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

// Обработка POST запроса удаления пользователя
// @Summary DeleteUser
// @Tags people
// @Description Удаление пользователя
// @ID deleteUser
// @Accept  json
// @Produce  json
// @Param passport body string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /deleteUser [post]
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

// Обработка POST запроса обновления пользователя
// @Summary UpdateUser
// @Tags people
// @Description Обновление пользователя
// @ID updateUser
// @Accept  json
// @Produce  json
// @Param params body classes.People true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /updateUser [post]
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

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.PUT_STRING, r.URL, params.String())
}

// Обработка POST запроса создания пользователя
// @Summary CreateUser
// @Tags people
// @Description Создание пользователя
// @ID createUser
// @Accept  json
// @Produce  json
// @Param passportSerie body string true "Параметры запроса"
// @Param passportNumber body string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /insertUser [post]
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

// Получение информации об одном пользователе
// @Summary UserInfo
// @Tags people
// @Description Получение информации об одном пользователе
// @ID userInfo
// @Accept  json
// @Produce  json
// @Param passportSerie query string true "Параметры запроса"
// @Param passportNumber query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /info [get]
func UserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Переменная для хранения параметров запроса
	passport := fmt.Sprintf("%s %s", r.URL.Query().Get("passportSerie"), r.URL.Query().Get("passportNumber"))

	// Вызов функции получения информации о пользователе
	db_data, err := db.GetUsers(db.GetUsersParams{Limit: 1}, classes.PeopleFilter{
		Passport: []string{passport},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to get user info %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(db_data[0].String()))
	
	log.Print(log.DEBUG, log.GET_STRING, r.URL, passport)
}

// Получение списка задач пользователя
// @Summary UserTasks
// @Tags people
// @Description Получение списка задач пользователя
// @ID userTasks
// @Accept  json
// @Produce  json
// @Param passportSerie query string true "Параметры запроса"
// @Param passportNumber query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /tasks [get]
func GetUserTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	passport := fmt.Sprintf("%s %s", r.URL.Query().Get("passportSerie"), r.URL.Query().Get("passportNumber"))
	log.Print(log.DEBUG, "Passport: %s", r.URL.Query())
	
	// Вызов функции получения информации о пользователе
	db_data, err := db.GetTasks(db.Psql, passport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to get user info %s", err.Error())
		return
	}

	// Преобразование полученных данных в JSON
	data, err := json.MarshalIndent(db_data, "", "    ")
	if err != nil {
		log.Print(log.ERROR, "Marshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
	
	log.Print(log.DEBUG, log.GET_STRING, r.URL, passport)
}

// Добавление задачи
// @Summary AddTask
// @Tags tasks
// @Description Добавление задачи
// @ID addTask
// @Accept  json
// @Produce  json
// @Param passport query string true "Параметры запроса"
// @Param title query string true "Параметры запроса"
// @Param description query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /addTask [post]
func AddTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Переменная для хранения параметров запроса
	var params struct {
		Passport string `json:"passport"`
		Title string `json:"title"`
		Description string `json:"description"`
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

	// Вызов функции добавления задачи
	err = db.AddTask(db.Psql, params.Title, params.Description, params.Passport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to add task %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.POST_STRING, r.URL, params.Passport, "OK")
}

// Начало задачи
// @Summary StartTask
// @Tags tasks
// @Description Начало задачи
// @ID startTask
// @Accept  json
// @Produce  json
// @Param taskId query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /startTask [put]
func StartTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Переменная для хранения параметров запроса
	taskId := r.URL.Query().Get("taskId")

	// Вызов функции начала задачи
	err := db.StartTaskTime(db.Psql, taskId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to start task %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.PUT_STRING, r.URL, taskId)
}

// Завершение задачи
// @Summary FinishTask
// @Tags tasks
// @Description Завершение задачи
// @ID finishTask
// @Accept  json
// @Produce  json
// @Param taskId query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /finishTask [put]
func FinishTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Переменная для хранения параметров запроса
	taskId := r.URL.Query().Get("taskId")

	// Вызов функции завершения задачи
	err := db.FinishTaskTime(db.Psql, taskId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to finish task %s", err.Error())
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	
	log.Print(log.DEBUG, log.PUT_STRING, r.URL, taskId)
}

// Получение времени
// @Summary GetUserTime
// @Tags tasks
// @Description Получение времени
// @ID getUserTime
// @Accept  json
// @Produce  json
// @Param passport query string true "Параметры запроса"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 404 {string} string "Not Found"
// @Router /userTime [get]
func GetUserTime(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Переменная для хранения параметров запроса
	passport := r.URL.Query().Get("passport")

	// Вызов функции получения времени
	time, tasksTimes, err := db.GetTime(db.Psql, passport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(log.ERROR, "Failed to get user time %s", err.Error())
		return
	}

	// Преобразование полученных данных в JSON
	data, err := json.MarshalIndent(map[string]any{
		"global": time,
		"tasks": tasksTimes,
	}, "", "    ")

	if err != nil {
		log.Print(log.ERROR, "Marshal error %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
	
	log.Print(log.DEBUG, log.GET_STRING, r.URL, passport)
}