package db

import (
	"database/sql"
	"fmt"
	"os"
	log "test_go_app/go/log"

	_ "github.com/lib/pq"
)

var Psql *sql.DB

func Connect() error {
	var err error
	pSqlConnect_String := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", ""),
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""))

	Psql, err = sql.Open("postgres", pSqlConnect_String)

	if err != nil {
		log.Print(log.ERROR, "Failed to connect to database")
		return fmt.Errorf("failed to connect to database")
	}
	// defer Psql.Close()

	err = Psql.Ping()
    if err != nil {
		log.Print(log.ERROR, "Failed to connect to database")
        return fmt.Errorf("failed to connect to database")
    }

	log.Print(log.INFO, "Connected to database")

	return nil
}

func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
		return value
    }

    return defaultVal
}