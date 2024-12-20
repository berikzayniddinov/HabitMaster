package main

import (
	"HabitMaster/databaseConnector"
	"HabitMaster/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Подключение к базе данных
	db := databaseConnector.ConnectDB()
	defer db.Close()

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Логирование запросов
	r.Use(loggingMiddleware)

	// Маршруты пользователей
	r.HandleFunc("/api/users", handlers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/api/get-users", handlers.GetUsers(db)).Methods("GET")

	// Маршруты привычек
	r.HandleFunc("/api/habits", handlers.CreateHabit(db)).Methods("POST")
	r.HandleFunc("/api/habits", handlers.GetHabits(db)).Methods("GET")
	r.HandleFunc("/api/habits", handlers.DeleteHabit(db)).Methods("DELETE")
	r.HandleFunc("/api/habits", handlers.UpdateHabit(db)).Methods("PUT")

	// Маршруты целей
	r.HandleFunc("/api/goals", handlers.CreateGoal(db)).Methods("POST")
	r.HandleFunc("/api/get-goals", handlers.GetGoals(db)).Methods("GET")
	r.HandleFunc("/api/goals", handlers.UpdateGoalByName(db)).Methods("PUT")
	r.HandleFunc("/api/goals/{id}", handlers.DeleteGoal(db)).Methods("DELETE")

	// Маршруты достижений
	r.HandleFunc("/api/achievements", handlers.CreateAchievement(db)).Methods("POST")
	r.HandleFunc("/api/achievements", handlers.GetAchievements(db)).Methods("GET")
	r.HandleFunc("/api/achievements/{id}", handlers.DeleteAchievementByID(db)).Methods("DELETE")

	// Маршруты прогресса
	r.HandleFunc("/api/progress", handlers.CreateProgress(db)).Methods("POST")
	r.HandleFunc("/api/progress", handlers.GetProgress(db)).Methods("GET")
	r.HandleFunc("/api/progress", handlers.UpdateProgressStatus(db)).Methods("PUT")

	// Маршруты уведомлений
	r.HandleFunc("/api/notifications", handlers.CreateNotification(db)).Methods("POST")
	r.HandleFunc("/api/notifications", handlers.GetNotifications(db)).Methods("GET")
	r.HandleFunc("/api/notifications/{id}", handlers.DeleteNotificationByID(db)).Methods("DELETE")

	// Раздача статических файлов
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./habittracker"))))

	// Запуск сервера
	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// loggingMiddleware — промежуточный слой для логирования входящих запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
