package handlers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var encryptionKey []byte

func CallbackGoogle(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	data, err := getUserData(state, code)
	if err != nil {
		log.Fatal("error getting user data")
	}
	userDataString := string(data) // Convertir los datos a un formato adecuado

	cookie := http.Cookie{
		Name:     "user_data",
		Value:    userDataString,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	fmt.Println(data)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "http://localhost:5173/", http.StatusTemporaryRedirect)
}
func getUserData(state, code string) ([]byte, error) {
	if state != RandomString {
		return nil, errors.New("invalid user state")
	}
	token, err := ssgolang.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
