package handlers

import (
	"encoding/json"
	"net/http"
	"backend/api/models"
	"backend/api/utils"
	"github.com/gorilla/mux"
	"backend/api/middleware"
	"github.com/golang-jwt/jwt/v5"

)

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error al decodificar los datos del usuario", http.StatusBadRequest)
		return
	}

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	createdUser := db.Create(&user)
	if createdUser.Error != nil {
		http.Error(w, "Error al insertar el usuario en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	db.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.MapClaims)
    if !ok {
        http.Error(w, "Información del usuario no disponible", http.StatusInternalServerError)
        return
    }

    uid := (*claims)["uid"].(string) 
	
	var user models.User

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", uid).First(&user)

	db.Where("id = ?", uid).First(&user)
	if user.UID == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", params["id"]).First(&user)

	if user.UID == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.Unscoped().Delete(&user)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario eliminado con éxito"))
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*jwt.MapClaims)
    if !ok {
        http.Error(w, "Información del usuario no disponible", http.StatusInternalServerError)
        return
    }

    uid := (*claims)["uid"].(string) 

	var user models.User

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	db.Where("id = ?", uid).First(&user)

	if user.UID == "" {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar solo los campos que deseas actualizar
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Error al decodificar los datos de actualización", http.StatusBadRequest)
		return
	}

	// Actualiza los campos específicos
	db.Model(&user).Updates(updates)

	// Envía una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Usuario actualizado con éxito"))
}
