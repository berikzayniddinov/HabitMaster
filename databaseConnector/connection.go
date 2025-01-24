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
	log.Println("Инициализация подключения к базе данных...")

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Ошибка загрузки файла окружения:", err)
	} else {
		log.Println("Файл окружения успешно загружен.")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Printf("Параметры подключения: host=%s, port=%s, user=%s, password=[скрыто], dbname=%s\n", dbHost, dbPort, dbUser, dbName)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)
	log.Println("Строка подключения сформирована:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	log.Println("Попытка проверить соединение с базой данных...")
	if err = db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v", err)
	}
	log.Println("Подключение к базе данных успешно установлено.")
	return db
}
