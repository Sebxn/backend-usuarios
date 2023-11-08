package models

type User struct {
	ID              int    `json:"ID" gorm:"column:id"`
	Nombre          string `json:"Nombre"`
	Apellido        string `json:"Apellido"`
	SegundoApellido string `json:"SegundoApellido" gorm:"column:segundoapellido"`
	Email           string `json:"Email"`
	Rut             string `json:"Rut"`
	Fono            string `json:"Fono"`
}

// Configura la tabla de la base de datos
func (User) TableName() string {
	return "usuario"
}
