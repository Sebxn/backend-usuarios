package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	
	"log"
	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v5"

	firebase "firebase.google.com/go"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // omitempty para no incluir el token si está vacío
	Error   string `json:"error,omitempty"` // omitempty para no incluir el error si está vacío
}

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

type FirebaseResponse struct {
	LocalId string `json:"localId"` // Assumiendo que el JSON de respuesta tiene un campo "localId"
}

func LoginUser(resp http.ResponseWriter, req *http.Request, app *firebase.App) {

	err := godotenv.Load(".env.credentials")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	// Define la URL del punto de conexión de Firebase Identity Toolkit para iniciar sesión con contraseña
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + os.Getenv("FIREBASE_API_KEY")// API KEY

	// Construye la carga útil de la solicitud JSON
	payload := fmt.Sprintf(`{
        "email": "%s",
        "password": "%s",
        "returnSecureToken": true
    }`, email, password)

	// Realiza la solicitud HTTP POST
	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))

	if err != nil {
		fmt.Fprintln(resp, "Error en la solicitud HTTP:", err)
		return
	}
	defer response.Body.Close()

	// Lee y procesa la respuesta
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(resp, "Error al leer la respuesta:", err)
		return
	}
	// Verifica el código de estado de la respuesta
	loginResp := LoginResponse{}

	if response.StatusCode == http.StatusOK {
		// Deserializa la respuesta JSON en la estructura FirebaseResponse.
		var fr FirebaseResponse
		if err := json.Unmarshal(responseBody, &fr); err != nil {
			fmt.Fprintln(resp, "Error al deserializar la respuesta:", err)
			return
		}
		// Generar un token JWT para el usuario registrado
		jwtToken, err := generateJWTToken(fr.LocalId)
		if err != nil {
			// Construir la respuesta de error
			loginResp.Error = "Error al generar el token JWT"
			loginResp.Status = "401"
			resp.WriteHeader(http.StatusInternalServerError)
		} else {
			// Construir la respuesta de éxito
			loginResp.Message = "Inicio de sesion correcto"
			loginResp.Token = jwtToken
			loginResp.Status = "200"
		}
	} else {
		// Autenticación fallida
		fmt.Fprintln(resp, "Error de autenticación:", string(responseBody))
	}
	resp.Header().Set("Content-Type", "application/json")

	// Enviar la respuesta en formato JSON
	json.NewEncoder(resp).Encode(loginResp)
	// Puedes procesar la respuesta de Firebase según tus necesidades

	// No es necesario devolver nada aquí, ya que la respuesta se maneja en el lugar.
}
