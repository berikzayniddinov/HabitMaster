package databaseConnector

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ConnectDB подключается к базе данных PostgreSQL и возвращает *sql.DB
func ConnectDB() *sql.DB {
	// Загрузка переменных окружения
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка загрузки файла окружения:", err)
	}

	// Чтение параметров подключения
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка проверки соединения с базой данных:", err)
	}

	fmt.Println("Подключение к базе данных успешно установлено.")
	return db
}
