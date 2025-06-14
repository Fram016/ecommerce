package models

import (
	"database/sql"
	"log"
)

// ProductoImagen representa la estructura de la tabla 'producto_imagenes'
type ProductoImagen struct {
	ID         int    `json:"id"`
	ProductoID int    `json:"producto_id"`
	RutaImagen string `json:"ruta_imagen"`
	TipoImagen string `json:"tipo_imagen"` // 'principal' o 'galeria'
}

// ListarImagenes devuelve todas las imágenes de un producto
func ListarImagenes(db *sql.DB, productoID int) ([]ProductoImagen, error) {
	rows, err := db.Query(`SELECT id, producto_id, ruta_imagen, tipo_imagen FROM producto_imagenes WHERE producto_id = ?`, productoID)
	if err != nil {
		log.Println("Error al obtener imágenes:", err)
		return nil, err
	}
	defer rows.Close()

	var imagenes []ProductoImagen
	for rows.Next() {
		var imagen ProductoImagen
		if err := rows.Scan(&imagen.ID, &imagen.ProductoID, &imagen.RutaImagen, &imagen.TipoImagen); err != nil {
			log.Println("Error al escanear imagen:", err)
			return nil, err
		}
		imagenes = append(imagenes, imagen)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return imagenes, nil
}

// CrearImagen inserta una nueva imagen para un producto
func CrearImagen(db *sql.DB, imagen ProductoImagen) error {
	query := `INSERT INTO producto_imagenes (producto_id, ruta_imagen, tipo_imagen) VALUES (?, ?, ?)`
	_, err := db.Exec(query, imagen.ProductoID, imagen.RutaImagen, imagen.TipoImagen)
	if err != nil {
		log.Println("Error al crear imagen:", err)
		return err
	}
	return nil
}

// EliminarImagen elimina una imagen de la base de datos
func EliminarImagen(db *sql.DB, id int) error {
	query := `DELETE FROM producto_imagenes WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar imagen:", err)
		return err
	}
	return nil
}
