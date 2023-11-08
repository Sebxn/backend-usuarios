package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func CallbackFacebook(w http.ResponseWriter, r *http.Request) {
	// Obtener el código de autorización de Facebook.
	code := r.URL.Query().Get("code")

	// Intercambiar el código de autorización por un token de acceso de Facebook.
	token, err := ssgolangf.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al intercambiar el código: %s", err), http.StatusInternalServerError)
		return
	}

	// Utilizar el token de acceso para hacer una solicitud a la API de Facebook.
	userID, userName, err := GetUserProfile(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener el perfil del usuario: %s", err), http.StatusInternalServerError)
		return
	}

	// Puedes hacer lo que desees con la información del usuario.
	fmt.Fprintf(w, "Token de acceso obtenido: %s\n", token.AccessToken)
	fmt.Fprintf(w, "ID de usuario: %s\n", userID)
	fmt.Fprintf(w, "Nombre de usuario: %s\n", userName)
}

func GetUserProfile(token *oauth2.Token) (string, string, error) {
	client := ssgolangf.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/v13.0/me?fields=id,name")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var userData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return "", "", err
	}

	return userData.ID, userData.Name, nil
}
