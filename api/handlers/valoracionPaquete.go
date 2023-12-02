package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"net/http"
)

// ValorarReserva maneja la valoración de una reserva.
func ValorarReserva(w http.ResponseWriter, r *http.Request) {
	var valoracion models.Valoracion

	err := json.NewDecoder(r.Body).Decode(&valoracion)
	if err != nil {
		http.Error(w, "Error al decodificar los datos de valoración", http.StatusBadRequest)
		return
	}

	// Verificar que la puntuación está en el rango correcto (1-5)
	if valoracion.Estrellas < 1 || valoracion.Estrellas > 5 {
		http.Error(w, "La puntuación debe estar en el rango de 1 a 5 estrellas", http.StatusBadRequest)
		return
	}

	// Obtener la reserva desde la base de datos
	reserva, err := getReservaFromDB(valoracion.IDReserva)
	if err != nil {
		http.Error(w, "Error al obtener la reserva desde la base de datos", http.StatusInternalServerError)
		return
	}

	// Verificar si la reserva está finalizada antes de permitir la valoración
	if reserva.Estado != "Finalizado" {
		http.Error(w, "No puedes valorar la reserva hasta que esté finalizada", http.StatusBadRequest)
		return
	}

	// Guardar la valoración en la base de datos
	if err := saveValoracionToDB(&valoracion); err != nil {
		http.Error(w, "Error al guardar la valoración de la reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Valoración realizada con éxito"})
}

// Función auxiliar para obtener una reserva desde la base de datos
// Función auxiliar para obtener una reserva desde la base de datos
func getReservaFromDB(reservaID int) (models.Reserva, error) {
	db, err := utils.OpenDBGorm()
	if err != nil {
		return models.Reserva{}, err
	}

	var reserva models.Reserva
	if err := db.Where("id = ?", reservaID).First(&reserva).Error; err != nil {
		return models.Reserva{}, err
	}

	return reserva, nil
}

// Función auxiliar para guardar una reserva en la base de datos
// Función para guardar una valoración en la base de datos
func saveValoracionToDB(valoracion *models.Valoracion) error {
	db, err := utils.OpenDBGorm()
	if err != nil {
		return err
	}

	// Guardar la valoración en la base de datos
	if err := db.Save(valoracion).Error; err != nil {
		return err
	}

	return nil
}
