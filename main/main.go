package main

import (
	"HabitMaster/databaseConnector"
	"HabitMaster/emailSender"
	"HabitMaster/handlers"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"sync"
)

var (
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

// getLimiter возвращает лимитер для IP клиента
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Если лимитер для IP не существует, создаем новый
	if limiter, exists := clients[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(1, 1) // 1 запрос в секунду, burst до 5
	clients[ip] = limiter
	return limiter
}

// rateLimiterMiddleware создает middleware для ограничения запросов
func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Получаем IP клиента
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware добавляет user_id в контекст из токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Пример декодирования токена
		userID, err := DecodeTokenAndGetUserID(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DecodeTokenAndGetUserID заглушка декодирования токена
func DecodeTokenAndGetUserID(token string) (int, error) {
	if token == "valid-token" {
		return 1, nil
	}
	return 0, fmt.Errorf("invalid token")
}

func main() {
	// Создаем лог-файл
	logFile, err := os.OpenFile("server_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Настройка логирования
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(logFile) // Перенаправляем вывод логов в файл

	log.Info("Инициализация сервера")

	// Подключение к базе данных
	db := databaseConnector.ConnectBD()
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Ошибка при закрытии соединения с базой данных")
		}
	}()

	log.Info("Успешное подключение к базе данных")

	// Создаем email-сервис
	emailService := emailSender.NewEmailSender()

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Логирование запросов
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(logrus.Fields{
				"method": r.Method,
				"url":    r.RequestURI,
				"remote": r.RemoteAddr,
			}).Info("Получен запрос")
			next.ServeHTTP(w, r)
		})
	})

	// Применяем middleware для ограничения запросов
	r.Use(rateLimiterMiddleware)

	// Определяем маршруты для профиля пользователя

	r.Handle("/api/profile", AuthMiddleware(http.HandlerFunc(handlers.GetUserProfile(db)))).Methods("GET")                    // Получение профиля
	r.Handle("/api/profile/update", AuthMiddleware(http.HandlerFunc(handlers.UpdateUserProfile(db)))).Methods("PATCH")        // Обновление профиля
	r.Handle("/api/profile/password", AuthMiddleware(http.HandlerFunc(handlers.ChangeUserPassword(db)))).Methods("POST")      // Смена пароля
	r.Handle("/api/profile/picture", AuthMiddleware(http.HandlerFunc(handlers.UploadUserProfilePicture(db)))).Methods("POST") // Загрузка фотографии профиля

	// Маршрут для главной страницы
	r.HandleFunc("/main.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./habittracker/main.html")
	}).Methods("GET")

	r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")
	r.Handle("/api/signup", http.HandlerFunc(handlers.CreateUser(db))).Methods("POST")  // Регистрация пользователя
	r.Handle("/api/login", http.HandlerFunc(handlers.LoginHandler(db))).Methods("POST") // Вход пользователя

	// Маршруты пользователей
	r.HandleFunc("/api/users", handlers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/api/get-users", handlers.GetUsers(db)).Methods("GET")

	// Маршруты привычек
	r.HandleFunc("/api/habits", handlers.CreateHabit(db)).Methods("POST")
	r.HandleFunc("/api/habits", handlers.GetHabits(db)).Methods("GET")
	r.HandleFunc("/api/habits", handlers.DeleteHabitByName(db)).Methods("DELETE")
	r.HandleFunc("/api/habits", handlers.UpdateHabit(db)).Methods("PUT")

	// Маршруты целей
	r.HandleFunc("/api/goals", handlers.CreateGoal(db)).Methods("POST")
	r.HandleFunc("/api/goals", handlers.GetGoals(db)).Methods("GET")
	r.HandleFunc("/api/goals", handlers.UpdateGoal(db)).Methods("PUT")
	r.HandleFunc("/api/goals", handlers.DeleteGoalByName(db)).Methods("DELETE")

	// Маршруты достижений
	r.HandleFunc("/api/achievements", handlers.CreateAchievement(db)).Methods("POST")
	r.HandleFunc("/api/achievements", handlers.GetAchievements(db)).Methods("GET")
	r.HandleFunc("/api/achievements", handlers.UpdateAchievement(db)).Methods("PUT")
	r.HandleFunc("/api/achievements", handlers.DeleteAchievementByTitle(db)).Methods("DELETE")

	// Маршруты прогресса
	r.HandleFunc("/api/progress", handlers.CreateProgress(db)).Methods("POST")
	r.HandleFunc("/api/progress", handlers.GetProgress(db)).Methods("GET")
	r.HandleFunc("/api/progress", handlers.UpdateProgressStatus(db)).Methods("PUT")

	// Маршруты уведомлений
	r.HandleFunc("/api/notifications", handlers.CreateNotification(db)).Methods("POST")
	r.HandleFunc("/api/notifications", handlers.GetNotifications(db)).Methods("GET")
	r.HandleFunc("/api/notifications", handlers.UpdateNotification(db)).Methods("PUT")
	r.HandleFunc("/api/notifications", handlers.DeleteNotificationByName(db)).Methods("DELETE")

	// Маршруты для отправки email
	r.HandleFunc("/api/admin/send-mass-email", handlers.SendMassEmailHandler(emailService)).Methods("POST")
	r.HandleFunc("/api/user/send-support-email", handlers.SendSupportEmailHandler(emailService)).Methods("POST")
	r.HandleFunc("/api/admin/send-email-with-attachment", handlers.SendEmailWithAttachmentHandler(emailService)).Methods("POST")

	// Раздача статических файлов
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./habittracker"))))

	// Запуск сервера
	log.Info("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.WithError(err).Fatal("Ошибка при запуске сервера")
	}
}
