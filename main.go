package main

import (
	"net/http"

	log "test_go_app/go/logLevel"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFound);
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowed);

	router.Handler("GET", "/", http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", router)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Print(log.WARNING, "Not Found: %s", r.URL)
	http.Error(w, "Not Found", http.StatusNotFound)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	log.Print(log.WARNING, "Method Not Allowed: %s", r.Method)
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

