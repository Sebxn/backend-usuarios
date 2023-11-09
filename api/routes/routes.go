package routes

import (
	"backend/api/handlers"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router, app *firebase.App) {
	// allowedOrigins := []string{"http://facturacion.lumonidy.studio", "http://localhost:3000"}

	// c := middleware.CorsMiddleware(allowedOrigins)
	// r.Use(c)
	r.Handle("/user/register", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(w, r, app)
	})).Methods("POST")

	r.Handle("/user/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, app)
	})).Methods("POST")

	r.Handle("/user/reset-password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ResetPassword(w, r)
	})).Methods("POST")

	r.Handle("/login-google", http.HandlerFunc(handlers.LoginGoogle))
	r.Handle("/login-facebook", http.HandlerFunc(handlers.LoginFacebook))

	r.HandleFunc("/users", handlers.AddUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserById).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PATCH")
	r.HandleFunc("/users/update_profile/{id}", handlers.UpdateProfile).Methods("PATCH")

	r.Handle("/api/facturacion/user_paquetes", http.HandlerFunc(handlers.ObtenerPaquetesByUser))
	r.Handle("/api/facturacion/actualizar_estado", http.HandlerFunc(handlers.ActualizarEstadoReserva))

	// Agrega más configuraciones de rutas aquí si es necesario
}
