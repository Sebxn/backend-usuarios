package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DBURL genera la URL de conexión a la base de datos PostgreSQL
func DBURL() string {
	err := godotenv.Load(".env.db")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DBHost := os.Getenv("DBHost")
	DBUser := os.Getenv("DBUser")
	DBPassword := os.Getenv("DBPassword")
	DBPort := os.Getenv("DBPort")
	DBName := os.Getenv("DBName")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName)
}
