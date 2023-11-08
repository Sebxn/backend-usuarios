package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var ssgolangf *oauth2.Config

func init() {
	err := godotenv.Load(".env.credentials")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ssgolangf = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID_FACEBOOK"),
		ClientSecret: os.Getenv("CLIENT_SECRET_FACEBOOK"),
		RedirectURL:  os.Getenv("REDIRECT_URL_FACEBOOK"),
		Endpoint:     facebook.Endpoint,
		Scopes:       []string{"public_profile", "email"},
	}
}

func LoginFacebook(w http.ResponseWriter, r *http.Request) {
	// Generar la URL de inicio de sesión de Facebook.
	url := ssgolangf.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}
