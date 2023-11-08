package handlers

import (
	"bytes"

	"fmt"
	"io/ioutil"
	"net/http"

	firebase "firebase.google.com/go"
)

func LoginUser(resp http.ResponseWriter, req *http.Request, app *firebase.App) {

	email := req.FormValue("email")
	password := req.FormValue("password")

	// Define la URL del punto de conexión de Firebase Identity Toolkit para iniciar sesión con contraseña
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyAN5HyOCaS1NMqiVKgcYaN1s6fq3oJWbMw" // API KEY

	// Construye la carga útil de la solicitud JSON
	payload := fmt.Sprintf(`{
        "email": "%s",
        "password": "%s",
        "returnSecureToken": true
    }`, email, password)

	// Realiza la solicitud HTTP POST
	response, err := http.Post(url, "application/login", bytes.NewBuffer([]byte(payload)))

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
	if response.StatusCode == http.StatusOK {
		// Autenticación exitosa
		fmt.Fprintln(resp, "Usuario autenticado con éxito.")
	} else {
		// Autenticación fallida
		fmt.Fprintln(resp, "Error de autenticación:", string(responseBody))
	}

	// Puedes procesar la respuesta de Firebase según tus necesidades

	// No es necesario devolver nada aquí, ya que la respuesta se maneja en el lugar.
}
