package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetReservasPorUsuario(w http.ResponseWriter, r *http.Request) {
	var reservas []models.Reserva

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	var requestBody map[string]string
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error al analizar el cuerpo JSON", http.StatusBadRequest)
		return
	}

	email, ok := requestBody["email"]
	if !ok {
		http.Error(w, "El campo 'email' es requerido en el cuerpo JSON", http.StatusBadRequest)
		return
	}

	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT * FROM reservas WHERE id_usuario = ?", email)
	if err != nil {
		http.Error(w, "Error al ejecutar la consulta SQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var reserva models.Reserva
		err := rows.Scan(&reserva.ID, &reserva.IDUsuario, &reserva.Estado, &reserva.FechaReserva)
		if err != nil {
			http.Error(w, "Error al escanear filas", http.StatusInternalServerError)
			return
		}
		reservas = append(reservas, reserva)
	}

	if len(reservas) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&reservas)
}

func GetReservasPorEstado(w http.ResponseWriter, r *http.Request) {
	var reservas []models.Reserva

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	var requestBody map[string]string
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error al analizar el cuerpo JSON", http.StatusBadRequest)
		return
	}

	email, ok := requestBody["email"]
	if !ok {
		http.Error(w, "El campo 'email' es requerido en el cuerpo JSON", http.StatusBadRequest)
		return
	}

	estado, ok := requestBody["estado"]
	if !ok {
		http.Error(w, "El campo 'estado' es requerido en el cuerpo JSON", http.StatusBadRequest)
		return
	}

	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT * FROM reservas WHERE id_usuario = ? AND estado = ?", email, estado)
	if err != nil {
		http.Error(w, "Error al ejecutar la consulta SQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var reserva models.Reserva
		err := rows.Scan(&reserva.ID, &reserva.IDUsuario, &reserva.Estado, &reserva.FechaReserva)
		if err != nil {
			http.Error(w, "Error al escanear filas", http.StatusInternalServerError)
			return
		}
		reservas = append(reservas, reserva)
	}

	if len(reservas) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&reservas)
}

func BorrarReserva(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	var requestBody map[string]interface{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error al analizar el cuerpo JSON", http.StatusBadRequest)
		return
	}

	id, ok := requestBody["id"].(float64)
	if !ok {
		http.Error(w, "El campo 'id' es requerido en el cuerpo JSON y debe ser un número", http.StatusBadRequest)
		return
	}

	idInt := int(id)

	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}

	var reserva models.Reserva
	err = db.QueryRow("SELECT * FROM reservas WHERE id = ?", idInt).
		Scan(&reserva.ID, &reserva.IDUsuario, &reserva.Estado, &reserva.FechaReserva)
	if err != nil {
		http.Error(w, "Reserva no encontrada", http.StatusNotFound)
		return
	}

	if reserva.Estado != "Aprobado" {
		http.Error(w, "No se puede borrar la reserva, no está en estado 'Aprobado'", http.StatusForbidden)
		return
	}

	diferenciaDias := time.Since(reserva.FechaReserva).Hours() / 24

	if diferenciaDias > 3 {
		http.Error(w, "No se puede borrar la reserva, han pasado más de 3 días desde la reserva", http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM reservas WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, "Error al borrar la reserva", http.StatusInternalServerError)
		return
	}

	mensaje := fmt.Sprintf("Reserva con ID %d borrada exitosamente", idInt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": mensaje})
}
