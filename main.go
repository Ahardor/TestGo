package main

import (
	"net/http"

	log "test_go_app/go/log"

	handlers "test_go_app/go/http"

	db "test_go_app/go/db"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Print(log.ERROR, "No .env file found")
		panic(err)
    }
}

func main() {
	if err := db.Connect(); err != nil {
		return 
	}

	// Запуск сервера
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFound);
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowed);

	router.GET("/", handlers.Index)
	router.POST("/getPeople", handlers.GetPeople)
	router.DELETE("/deleteUser", handlers.DeleteUser)
	router.PUT("/updateUser", handlers.UpdateUser)
	router.POST("/insertUser", handlers.CreateUser)

	log.Print(log.INFO, "Server starting on port 8080")
	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Print(log.ERROR, "Failed to start server")
	}	
}

// Обработчики ошибок
func notFound(w http.ResponseWriter, r *http.Request) {
	log.Print(log.WARNING, "Not Found: %s", r.URL)
	http.Error(w, "Not Found", http.StatusNotFound)
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	log.Print(log.WARNING, "Method Not Allowed: %s", r.Method)
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}