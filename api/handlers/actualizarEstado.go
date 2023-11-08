package handlers

import (
	"backend/api/utils"
	"net/http"
)

func ActualizarEstadoReserva(w http.ResponseWriter, r *http.Request) {
	// Obtener datos de la solicitud, como el ID de reserva y el nuevo estado
	reservaID := r.FormValue("reservaID")
	nuevoEstado := r.FormValue("nuevoEstado")

	// Realizar la actualización en la base de datos (usando SQL, por ejemplo)
	// Supongamos que tienes una función "actualizarEstadoReservaEnBD" que actualiza el estado en la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al abrir la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza una consulta para verificar si la reserva con el ID especificado existe
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM reserva WHERE id = $1)", reservaID).Scan(&exists)
	if err != nil {
		http.Error(w, "Error al verificar la existencia de la reserva", http.StatusInternalServerError)
		return
	}

	// Si la reserva no existe, devuelve una respuesta de error
	if !exists {
		http.Error(w, "La reserva con el ID especificado no existe", http.StatusNotFound)
		return
	}

	// Define la consulta SQL para actualizar el estado de la reserva
	query := "UPDATE reserva SET estado = $1 WHERE id = $2"

	// Ejecuta la consulta SQL
	_, err = db.Exec(query, nuevoEstado, reservaID)
	if err != nil {
		http.Error(w, "Error al actualizar el estado de la reserva", http.StatusInternalServerError)
		return
	}

	// Enviar una respuesta de éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Estado de la reserva actualizado correctamente"))
}
