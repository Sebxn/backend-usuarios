package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	apiKey := "AIzaSyAN5HyOCaS1NMqiVKgcYaN1s6fq3oJWbMw" // Reemplaza con tu API Key
	resetPasswordURL := "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=" + apiKey
	payload := `{
        "requestType": "PASSWORD_RESET",
        "email": "` + email + `"
    }`

	response, err := http.Post(resetPasswordURL, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		// La solicitud no se completó con éxito; maneja el error aquí
		fmt.Printf("Error al solicitar restablecimiento de contraseña. Código de estado: %d\n", response.StatusCode)
		fmt.Printf("Respuesta: %s\n", string(body))
		return nil
	}

	// La solicitud se completó con éxito; el correo de restablecimiento de contraseña se envió
	fmt.Println("Correo de restablecimiento de contraseña enviado con éxito")
	return nil
}
