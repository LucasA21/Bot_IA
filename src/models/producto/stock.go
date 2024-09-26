package stock

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type Producto struct {
	UserID   int64  `json:"user_id"`
	Producto string `json:"producto"`
	Cantidad int    `json:"cantidad"`
}

func GenerarArchivoStock(conn *pgx.Conn, userID int64) (string, error) {
	// Consulta el stock del usuario
	rows, err := conn.Query(context.Background(),
		"SELECT producto, cantidad FROM stock WHERE user_id = $1", userID)
	if err != nil {
		return "", fmt.Errorf("error al consultar el stock: %v", err)
	}
	defer rows.Close()

	// Crear archivo temporal (sobrescribiendo si ya existe)
	fileName := fmt.Sprintf("stock_%d.txt", userID)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo: %v", err)
	}
	defer file.Close()

	// Escribir el stock en el archivo
	for rows.Next() {
		var producto string
		var cantidad int
		err := rows.Scan(&producto, &cantidad)
		if err != nil {
			return "", fmt.Errorf("error al leer fila: %v", err)
		}
		file.WriteString(fmt.Sprintf("%s: %d\n", producto, cantidad))
	}

	// Verificar si hubo errores en la iteración de filas
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error en la iteración de filas: %v", err)
	}

	return fileName, nil
}

func AgregarAlStock(conn *pgx.Conn, userID int64, producto string, cantidad int) error {
	var currentQuantity int
	err := conn.QueryRow(context.Background(),
		"SELECT cantidad FROM stock WHERE user_id = $1 AND producto = $2", userID, producto).Scan(&currentQuantity)

	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = conn.Exec(context.Background(),
				"INSERT INTO stock (user_id, producto, cantidad) VALUES ($1, $2, $3)", userID, producto, cantidad)
			if err != nil {
				log.Printf("Error al agregar nuevo producto: %v", err)
				return fmt.Errorf("error al agregar al stock: %v", err)
			}
			log.Printf("Producto '%s' agregado al stock por el usuario %d.", producto, userID)
		} else {
			log.Printf("Error al consultar el stock: %v", err)
			return fmt.Errorf("error al consultar el stock: %v", err)
		}
	} else {
		newQuantity := currentQuantity + cantidad
		_, err = conn.Exec(context.Background(),
			"UPDATE stock SET cantidad = $1 WHERE user_id = $2 AND producto = $3", newQuantity, userID, producto)
		if err != nil {
			log.Printf("Error al actualizar el producto: %v", err)
			return fmt.Errorf("error al actualizar el stock: %v", err)
		}
		log.Printf("Producto '%s' actualizado. Nueva cantidad: %d para el usuario %d.", producto, newQuantity, userID)
	}
	return nil
}
