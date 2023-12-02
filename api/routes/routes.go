package routes

import (
	"backend/api/handlers"
	"backend/api/middleware"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

func ConfigureRoutes(r *mux.Router, app *firebase.App) {
	allowedOrigins := []string{"http://localhost:5173"} // Replace with your allowed origins
	r.Use(middleware.CorsMiddleware(allowedOrigins))

	// r.Handle("/prueba", middleware.TokenAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.Prueba(w, r)
	// }))).Methods("GET")

	r.Handle("/user/register", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterUser(w, r, app)
	})).Methods("POST")
	r.Handle("/user/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(w, r, app)
	})).Methods("POST")
	r.Handle("/user/reset-password", middleware.TokenAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ResetPassword(w, r)
	}))).Methods("POST")

	r.Handle("/user", middleware.TokenAuthMiddleware(http.HandlerFunc(handlers.GetUserById))).Methods("GET", "OPTIONS")
	r.Handle("/user/update", middleware.TokenAuthMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PATCH", "OPTIONS")
	//paquetes comprados por el usuario
	r.Handle("/user/mis_reservas", middleware.TokenAuthMiddleware(http.HandlerFunc(handlers.ObtenerPaquetesByUser))).Methods("GET", "OPTIONS")

	r.HandleFunc("/users/update_profile/{id}", handlers.UpdateProfile).Methods("PATCH") // ???

	r.Handle("/login-google", http.HandlerFunc(handlers.LoginGoogle))
	r.Handle("/login-facebook", http.HandlerFunc(handlers.LoginFacebook))

	r.HandleFunc("/users", handlers.AddUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	r.Handle("/api/facturacion/actualizar_estado", http.HandlerFunc(handlers.ActualizarEstadoReserva))

	r.HandleFunc("/reservas/valorar", handlers.ValorarReserva).Methods("POST")

	// Agrega más configuraciones de rutas aquí si es necesario
}
