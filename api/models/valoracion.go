package models

type Valoracion struct {
	ID         int    `json:"id"`
	IDReserva  int    `json:"id_reserva"`
	Estrellas  int    `json:"estrellas"`
	Comentario string `json:"comentario"`
}

func (Valoracion) TableName() string {
	return "valoracion"
}
