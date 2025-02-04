package main

import (
	"HabitMaster/auth"
	"HabitMaster/databaseConnector"
	"HabitMaster/emailSender"
	"HabitMaster/handlers"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strings"
	"sync"
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

var (
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := clients[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(1, 1)
	clients[ip] = limiter
	return limiter
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		userID, err := DecodeTokenAndGetUserID(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeTokenAndGetUserID(authHeader string) (int, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("invalid authorization header")
	}

	tokenString := parts[1]
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

func main() {
	logFile, err := os.OpenFile("server_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(logFile)

	log.Info("Инициализация сервера")

	db := databaseConnector.ConnectBD()
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Ошибка при закрытии соединения с базой данных")
		}
	}()

	log.Info("Успешное подключение к базе данных")

	emailService := emailSender.NewEmailSender()
	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.WithFields(logrus.Fields{
				"method": req.Method,
				"url":    req.RequestURI,
				"remote": req.RemoteAddr,
			}).Info("Получен запрос")
			next.ServeHTTP(w, req)
		})
	})

	r.Use(rateLimiterMiddleware)

	// Пример защищённого роутера
	protected := r.PathPrefix("/api/protected").Subrouter()
	protected.Use(AuthMiddleware)

	protected.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(int)
		fmt.Fprintf(w, "This is a protected route. user_id = %d\n", userID)
	}).Methods("GET")

	// Фронтенд
	r.HandleFunc("/main.html", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./habittracker/main.html")
	}).Methods("GET")

	r.HandleFunc("/register", auth.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify-email", auth.VerifyCode).Methods(http.MethodPost)
	r.HandleFunc("/logout", auth.Logout).Methods(http.MethodPost)

	// Привычки
	r.HandleFunc("/api/habits", handlers.CreateHabit(db)).Methods("POST")
	r.HandleFunc("/api/habits", handlers.GetHabits(db)).Methods("GET")
	r.HandleFunc("/api/habits", handlers.DeleteHabitByName(db)).Methods("DELETE")
	r.HandleFunc("/api/habits", handlers.UpdateHabit(db)).Methods("PUT")

	// Роли и авторизация
	r.HandleFunc("/api/assign-role", handlers.AssignRoleToUser(db)).Methods("POST")
	r.Handle("/api/admin-action", handlers.RoleMiddleware("admin", db)(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("This is an admin action."))
	})))

	// Цели
	r.HandleFunc("/api/goals", handlers.CreateGoal(db)).Methods("POST")
	r.HandleFunc("/api/goals", handlers.GetGoals(db)).Methods("GET")
	r.HandleFunc("/api/goals", handlers.UpdateGoal(db)).Methods("PUT")
	r.HandleFunc("/api/goals", handlers.DeleteGoalByName(db)).Methods("DELETE")
	r.HandleFunc("/api/goals/deleteAll", handlers.DeleteAllGoals(db)).Methods("DELETE") // Новый маршрут

	// Email-уведомления
	r.HandleFunc("/api/admin/send-mass-email", handlers.SendMassEmailHandler(emailService)).Methods("POST")
	r.HandleFunc("/api/user/send-support-email", handlers.SendSupportEmailHandler(emailService)).Methods("POST")
	r.HandleFunc("/api/admin/send-email-with-attachment", handlers.SendEmailWithAttachmentHandler(emailService)).Methods("POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./habittracker"))))

	log.Info("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.WithError(err).Fatal("Ошибка при запуске сервера")
	}
}
