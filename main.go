package main

// @title           Test APP GO
// @version         1.0
// @description     Тестовое задание по GO

// @host      localhost:8080
// @BasePath  /

import (
	"net/http"

	log "test_go_app/pkg/log"

	handlers "test_go_app/pkg/http"

	db "test_go_app/pkg/db"

	_ "test_go_app/docs"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
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

	router.ServeFiles("/docs/*filepath", http.Dir("./docs"))

	router.GET("/", handlers.Index)
	router.GET("/info", handlers.UserInfo)
	router.POST("/getPeople", handlers.GetPeople)
	router.DELETE("/deleteUser", handlers.DeleteUser)
	router.PUT("/updateUser", handlers.UpdateUser)
	router.POST("/insertUser", handlers.CreateUser)
	router.GET("/tasks", handlers.GetUserTasks)
	router.PUT("/startTask", handlers.StartTask)
	router.PUT("/finishTask", handlers.FinishTask)
	router.GET("/userTime", handlers.GetUserTime)
	router.POST("/addTask", handlers.AddTask)
	router.GET("/api/:any", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Print(log.INFO, "Swagger")
		httpSwagger.WrapHandler(w, r)
	})

	log.Print(log.INFO, "Server starting on port 8080")

	server := &http.Server{Addr: ":8080", Handler: router}
	if err := server.ListenAndServe(); err != nil {
		log.Print(log.ERROR, "Failed to start server")
	}

	server.ListenAndServe()
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