package databaseConnector

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectBD() *sql.DB {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка загрузки файла окружения:", err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка проверки соединения с базой данных:", err)
	}

	fmt.Println("Подключение к базе данных успешно установлено.")
	return db
}
