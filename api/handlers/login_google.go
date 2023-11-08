package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssgolang *oauth2.Config
var RandomString = "randoms-tring"

func init() {

	err := godotenv.Load(".env.credentials")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ssgolang = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL_GOOGLE"),
		ClientID:     os.Getenv("CLIENT_ID_GOOGLE"),
		ClientSecret: os.Getenv("CLIENT_SECRET_GOOGLE"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	url := ssgolang.AuthCodeURL(RandomString)
	// En lugar de redirigir, devolvemos la URL como una respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
