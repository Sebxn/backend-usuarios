package handlers

import (
	"encoding/json"

	"net/http"

	"backend/api/models"
	"backend/api/utils"
)

func ObtenerPaquetesByUser(w http.ResponseWriter, r *http.Request) {
	// Captura los parámetros de la URL
	id_usuario := r.URL.Query().Get("id_usuario")

	// Abre la conexión con la base de datos
	db, err := utils.OpenDB()
	if err != nil {
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Realiza la consulta SQL con los parámetros
	rows, err := db.Query(`
	WITH ranked_reservas AS (
		SELECT
			reserva.id,
			reserva.id_usuario,
			paquete.nombre,
			COALESCE(total_personas, 0) AS total_personas,
			fechapaquete.fechainit,
			fechapaquete.fechafin,
			ciudad_origen.nombre AS nombre_ciudad_origen,
			ciudad_destino.nombre AS nombre_ciudad_destino,
			fechapaquete.precio_oferta_vuelo as oferta_vuelo,
			paquete.precio_vuelo,
			habitacionhotel.precio_noche,
			paquete.descripcion,
			paquete.detalles,
			opcionhotel.nombre AS nombre_opcion_hotel,
			habitacionhotel.descripcion AS descripcion_habitacion,
			habitacionhotel.servicios AS servicios_habitacion,
			hotel.nombre AS nombre_hotel,
			hotel.direccion AS direccion_hotel,
			hotel.valoracion AS valoracion_hotel,
			hotel.descripcion AS descripcion_hotel,
			hotel.servicios AS servicios_hotel,
			hotel.telefono AS telefono_hotel,
			hotel.correo_electronico AS correo_electronico_hotel,
			hotel.sitio_web AS sitio_web_hotel,
			ROW_NUMBER() OVER (PARTITION BY reserva.id ORDER BY t.ord) AS row_num
		FROM
			reserva
			INNER JOIN fechapaquete ON fechapaquete.id = reserva.id_fechapaquete
			INNER JOIN paquete ON paquete.id = fechapaquete.id_paquete
			INNER JOIN ciudad ciudad_origen ON paquete.id_origen = ciudad_origen.id
			INNER JOIN ciudad ciudad_destino ON paquete.id_destino = ciudad_destino.id
			INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
			INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
			INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
			INNER JOIN hotel ON habitacionhotel.hotel_id = hotel.id
			LEFT JOIN (
				SELECT
					paquete.id AS paquete_id,
					SUM(opcionhotel.cantidad) AS total_personas
				FROM
					paquete
					INNER JOIN unnest(paquete.id_hh) WITH ORDINALITY t(habitacion_id, ord) ON TRUE
					INNER JOIN habitacionhotel ON t.habitacion_id = habitacionhotel.id
					INNER JOIN opcionhotel ON habitacionhotel.opcion_hotel_id = opcionhotel.id
				GROUP BY
					paquete.id
			) AS subquery ON paquete.id = subquery.paquete_id
	)
	SELECT *
	FROM ranked_reservas
	WHERE row_num = 1 AND id_usuario = $1
	ORDER BY id;
	`, id_usuario)
	if err != nil {
		http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Procesa los resultados y crea una estructura para la respuesta JSON
	var reservas []models.Reserva
	for rows.Next() {
		var reserva models.Reserva
		if err := rows.Scan(
			&reserva.ID,
			&reserva.IDUsuario, // Add IDUsuario here
			&reserva.Nombre,
			&reserva.TotalPersonas, // Add TotalPersonas here
			&reserva.FechaInicio,
			&reserva.FechaFin,
			&reserva.NombreCiudadOrigen,
			&reserva.NombreCiudadDestino,
			&reserva.OfertaVuelo, // Add OfertaVuelo here
			&reserva.PrecioVuelo,
			&reserva.PrecioNoche,
			&reserva.Descripcion,
			&reserva.Detalles,
			&reserva.NombreOpcionHotel, // Add NombreOpcionHotel here
			&reserva.DescripcionHabitacion,
			&reserva.ServiciosHabitacion,
			&reserva.NombreHotel,
			&reserva.DireccionHotel,
			&reserva.ValoracionHotel,
			&reserva.DescripcionHotel,
			&reserva.ServiciosHotel,
			&reserva.TelefonoHotel,
			&reserva.CorreoElectronicoHotel,
			&reserva.SitioWebHotel,
			&reserva.RowNum,
		); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}

		reservas = append(reservas, reserva)
	}
	// Convierte los resultados a JSON y envía la respuesta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reservas); err != nil {
		http.Error(w, "Error al convertir a JSON", http.StatusInternalServerError)
	}

	if reservas == nil {
		http.Error(w, "El usuario no existe", http.StatusInternalServerError)
		return
	}
}
