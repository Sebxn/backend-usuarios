package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt/v5"
)

func generateJWTToken(uid string) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims (payload)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	// Generate encoded token and send it as response
	jwtToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request, app *firebase.App) error {

	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
		return err
	}

	// Recuperar los datos del formulario
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	// Verificar si las contraseñas coinciden
	if password != confirmPassword {
		http.Error(w, "Las contraseñas no coinciden", http.StatusBadRequest)
		return nil
	}

	// Verificar la longitud y los requisitos de la contraseña
	if len(password) < 8 || len(password) > 32 {
		http.Error(w, "La contraseña debe tener entre 8 y 32 caracteres.", http.StatusBadRequest)
		return nil
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, char := range password {
		if 'A' <= char && char <= 'Z' {
			hasUpper = true
		}
		if 'a' <= char && char <= 'z' {
			hasLower = true
		}
		if '0' <= char && char <= '9' {
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		http.Error(w, "La contraseña debe contener al menos una mayúscula, una minúscula y un número.", http.StatusBadRequest)
		return nil
	}

	// Register the user with email and password
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)

	user, err := client.CreateUser(ctx, params)
	if err != nil {
		// Manejar errores
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return err
	}

	userPostgreSQL := models.User{
		UID:             user.UID,              // El UID del usuario de Firebase
		Nombre:          r.FormValue("Nombre"), // Recupera los valores del formulario
		Apellido:        r.FormValue("Apellido"),
		SegundoApellido: r.FormValue("SegundoApellido"),
		Email:           email,
		Rut:             r.FormValue("Rut"),
		Fono:            r.FormValue("Fono"),
	}

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return err
	}

	createdUser := db.Create(&userPostgreSQL)
	if createdUser.Error != nil {
		http.Error(w, "Error al insertar el usuario en la base de datos", http.StatusInternalServerError)
		return createdUser.Error
	}

	// Generar un token JWT para el usuario registrado
	jwtToken, err := generateJWTToken(user.UID)
	if err != nil {
		// Manejar errores
		http.Error(w, "Error al generar el token JWT", http.StatusInternalServerError)
		return err
	}

	// Respond with "Contraseña exitosa" y el JWT token
	fmt.Fprintln(w, "Contraseña exitosa")
	fmt.Fprintln(w, jwtToken)
	return nil
}
