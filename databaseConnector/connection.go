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
	log.Println("📌 Инициализация подключения к базе данных...")

	// Получаем текущую директорию
	cwd, _ := os.Getwd()
	log.Println("📂 Текущая директория:", cwd)

	// Проверяем наличие .env в текущей директории
	paths := []string{
		"C:/Users/berik/GolandProjects/HabitMaster/.env", // Абсолютный путь
		".env", "../.env", "../../.env", "Testing/.env",
	}
	envLoaded := false

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				log.Println("✅ Файл окружения загружен из:", path)
				envLoaded = true
				break
			} else {
				log.Println("❌ Ошибка загрузки:", path, err)
			}
		} else {
			log.Println("❌ Файл .env не найден в:", path)
		}
	}

	if !envLoaded {
		log.Fatal("❌ Файл .env не найден! Проверь расположение.")
	}

	// Проверяем переменные
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Printf("🔍 Загруженные переменные: host=%s, port=%s, user=%s, dbname=%s\n", dbHost, dbPort, dbUser, dbName)

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("❌ Переменные окружения не загружены!")
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Println("🔌 Строка подключения сформирована:", connStr)

	// Подключение к БД
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения: %v", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Ошибка проверки соединения: %v", err)
	}

	log.Println("✅ Подключение к базе данных успешно установлено.")
	return db
}
