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

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!")
	log.Print(log.INFO, log.GET_STRING, r.URL, "Hello World!")
}

func People(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var params db.GetUsersParams
	var filter classes.PeopleFilter
	
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Print(log.ERROR, "Read body error %s", err.Error())
		return
	}

    err = json.Unmarshal(data, &params)
	if err != nil {
		log.Print(log.ERROR, "Unmarshal error %s", err.Error())
		return
	}
	err = json.Unmarshal(data, &filter)
	if err != nil {
		log.Print(log.ERROR, "Unmarshal error %s", err.Error())
		return
	}

	params.Psql = db.Psql

	s, err := db.GetUsers(
		params,
		filter,
	)
	if err != nil {
		return
	}

	fmt.Fprint(w, s)
	log.Print(log.DEBUG, log.POST_STRING, r.URL, fmt.Sprintf("%s %s", filter.String(), params.String()), s)
}