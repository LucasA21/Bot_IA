package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

// Función para establecer conexión con la BDD
func ConnectToDB() *pgx.Conn {
	log.Printf("Conectando a la base de datos: %s", os.Getenv("DATABASE_URL"))

	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("No se pudo parsear la URL de la base de datos: %v", err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	return conn
}
