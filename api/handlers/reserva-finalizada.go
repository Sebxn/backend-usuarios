package handlers

import (
	"backend/api/models"
	"backend/api/utils"
	"net/http"
	"time"
)

func ReservaFinalizada(w http.ResponseWriter, r *http.Request) {
	var reserva models.Reserva
	fechaActual := time.Now()

	// Comprobar si la fecha actual es después de la fecha de terminación de la reserva
	if fechaActual.After(reserva.FechaFin) {
		// Actualizar el estado de la reserva a "Finalizado"
		reserva.Estado = "Finalizado"

		db, err := utils.OpenDBGorm()
		if err != nil {
			http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
			return
		}
		// Actualizar la reserva en la base de datos
		if err := db.Model(&reserva).Updates(models.Reserva{Estado: "Finalizado"}).Error; err != nil {
			// Manejar el error, por ejemplo, imprimirlo o devolverlo
			panic(err)
		}
	}
}
