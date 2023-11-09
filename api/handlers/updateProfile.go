package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User

	db, err := utils.OpenDBGorm()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	// Busca el usuario por ID
	db.First(&user, params["id"])

	if user.UID == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuario no encontrado"))
		return
	}

	// Decodifica el JSON de la solicitud
	var requestPayload map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Actualiza los campos seleccionados del usuario
	if val, ok := requestPayload["Nombre"]; ok {
		user.Nombre = val.(string)
	}
	if val, ok := requestPayload["Apellido"]; ok {
		user.Apellido = val.(string)
	}
	if val, ok := requestPayload["SegundoApellido"]; ok {
		user.SegundoApellido = val.(string)
	}

	if val, ok := requestPayload["Fono"]; ok {
		user.Fono = val.(string)
	}
	if val, ok := requestPayload["FotoPerfil"]; ok {
		user.FotoPerfil = val.(string)
	}

	// Actualiza el usuario en la base de datos
	db.Save(&user)

	// Responde con el usuario actualizado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}
