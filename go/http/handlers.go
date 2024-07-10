package handlers

import (
	"fmt"
	"net/http"
	log "test_go_app/go/log"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!")
	log.Print(log.INFO, log.GET_STRING, r.URL, "Hello World!")
}