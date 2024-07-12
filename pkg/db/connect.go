package db

import (
	"database/sql"
	"fmt"
	"os"
	log "test_go_app/pkg/log"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Глобальная переменная для соединения с базой данных
var Psql *sql.DB

// Подключение к базе данных
func Connect() error {

	// Выполнение миграции
	err := migrateDB_Down()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	err = migrateDB()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	// Строка подключения к базе данных с использованием переменных окружения
	pSqlConnect_String := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", ""),
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""))

	// Подключение к базе данных
	Psql, err = sql.Open("postgres", pSqlConnect_String)

	if err != nil {
		log.Print(log.ERROR, "Failed to connect to database")
		return fmt.Errorf("failed to connect to database")
	}

	// Проверка подключения
	err = Psql.Ping()
    if err != nil {
		log.Print(log.ERROR, "Failed to connect to database")
        return fmt.Errorf("failed to connect to database")
    }

	// Подключение к базе данных успешно
	log.Print(log.INFO, "Connected to database")

	return nil
}

// Получение переменной окружения
func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
		return value
    }

    return defaultVal
}

func migrateDB() error {
	postgresString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""),
		getEnv("DB_NAME", ""))

	m, err := migrate.New(
		"file://db/migrations",
		postgresString)
	if err != nil {
		log.Print(log.ERROR, err.Error())
		return err
	}
	if err := m.Up(); err != nil {
		log.Print(log.ERROR, err.Error())
		return err
	}

	return nil
}

func migrateDB_Down() error {
	postgresString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", ""),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", ""),
		getEnv("DB_PORT", ""),
		getEnv("DB_NAME", ""))

	m, err := migrate.New(
		"file://db/migrations",
		postgresString)
	if err != nil {
		log.Print(log.ERROR, err.Error())
		return err
	}
	if err := m.Down(); err != nil {
		log.Print(log.ERROR, err.Error())
		return err
	}

	return nil
}