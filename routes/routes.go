package routes

import (
	"HabitMaster/auth"

	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", auth.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", auth.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify-email", auth.VerifyCode).Methods(http.MethodPost)
	r.HandleFunc("/logout", auth.Logout).Methods(http.MethodPost)

	return r
}
