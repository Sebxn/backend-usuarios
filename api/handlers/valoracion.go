package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"encoding/json"
	"net/http"
)

func ValorarReserva(w http.ResponseWriter, r *http.Request) {
	var reserva models.Reserva
	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, "Error al decodificar los datos de valoración", http.StatusBadRequest)
		return
	}

	// Verificar que la puntuación está en el rango correcto (1-5)
	if reserva.Estrellas < 1 || reserva.Estrellas > 5 {
		http.Error(w, "La puntuación debe estar en el rango de 1 a 5 estrellas", http.StatusBadRequest)
		return
	}

	// Aquí deberías implementar la lógica para manejar la valoración en la base de datos
	err = saveValoracionToDB(&reserva)
	if err != nil {
		http.Error(w, "Error al guardar la valoración de la reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Valoración realizada con éxito"})
}

// Función auxiliar para guardar una valoración en la base de datos
func saveValoracionToDB(reserva *models.Reserva) error {
	db, err := utils.OpenDBGorm()
	if err != nil {
		return err
	}

	// Obtener la reserva desde la base de datos
	var reservaDB models.Reserva
	if err := db.Where("id = ?", reserva.ID).First(&reservaDB).Error; err != nil {
		return err
	}

	// Actualizar la puntuación del usuario en la reserva
	reservaDB.Estrellas = reserva.Estrellas

	// Guardar la reserva actualizada en la base de datos
	if err := db.Save(&reservaDB).Error; err != nil {
		return err
	}

	return nil
}
