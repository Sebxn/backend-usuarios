package utils

import (
	"backend/api/config"
	"database/sql"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// OpenDB abre la conexión con la base de datos y la devuelve como una instancia de GORM
func OpenDBGorm() (*gorm.DB, error) {
	dsn := config.DBURL() // Asegúrate de configurar tu DSN según tus necesidades

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DBURL())
	if err != nil {
		return nil, err
	}
	return db, nil
}
