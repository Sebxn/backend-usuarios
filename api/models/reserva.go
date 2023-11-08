package models

import (
	"time"
)

type Reserva struct {
	ID                     int       `json:"id"`
	IDUsuario              string    `json:"id_usuario"`
	Nombre                 string    `json:"nombre"`
	TotalPersonas          int       `json:"total_personas"`
	FechaInicio            time.Time `json:"fechainit"`
	FechaFin               time.Time `json:"fechafin"`
	NombreCiudadOrigen     string    `json:"nombre_ciudad_origen"`
	NombreCiudadDestino    string    `json:"nombre_ciudad_destino"`
	OfertaVuelo            float64   `json:"oferta_vuelo"`
	PrecioVuelo            float64   `json:"precio_vuelo"`
	PrecioNoche            float64   `json:"precio_noche"`
	Descripcion            string    `json:"descripcion"`
	Detalles               string    `json:"detalles"`
	NombreOpcionHotel      string    `json:"nombre_opcion_hotel"`
	DescripcionHabitacion  string    `json:"descripcion_habitacion"`
	ServiciosHabitacion    string    `json:"servicios_habitacion"`
	NombreHotel            string    `json:"nombre_hotel"`
	DireccionHotel         string    `json:"direccion_hotel"`
	ValoracionHotel        float64   `json:"valoracion_hotel"`
	DescripcionHotel       string    `json:"descripcion_hotel"`
	ServiciosHotel         string    `json:"servicios_hotel"`
	TelefonoHotel          string    `json:"telefono_hotel"`
	CorreoElectronicoHotel string    `json:"correo_electronico_hotel"`
	SitioWebHotel          string    `json:"sitio_web_hotel"`
	RowNum                 int       `json:"row_num"`
}
