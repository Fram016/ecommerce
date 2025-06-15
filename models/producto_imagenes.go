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

// ImagenPrincipal devuelve la imagen principal de un producto
func ImagenPrincipal(db *sql.DB, productoID int) (ProductoImagen, error) {
    row := db.QueryRow(`SELECT id, producto_id, ruta_imagen, tipo_imagen FROM producto_imagenes WHERE producto_id = ? AND tipo_imagen = "principal" LIMIT 1`, productoID)

    imagen := ProductoImagen{}
    if err := row.Scan(&imagen.ID, &imagen.ProductoID, &imagen.RutaImagen, &imagen.TipoImagen); err != nil {
        if err == sql.ErrNoRows {
            return ProductoImagen{}, nil // No se encuentra la imagen, devuelve un objeto vacío
        }
        log.Println("Error al escanear imagen:", err)
        return ProductoImagen{}, err // Retorna un error si falla la consulta
    }
    return imagen, nil // Devuelve la imagen completa si no hay error
}

//buscar imagen por id
func ObtenerImagen(db *sql.DB, id int) (ProductoImagen, error) {
	query := `SELECT id, producto_id, ruta_imagen, tipo_imagen FROM producto_imagenes WHERE id = ?`
	var imagen ProductoImagen
	err := db.QueryRow(query, id).Scan(&imagen.ID, &imagen.ProductoID, &imagen.RutaImagen, &imagen.TipoImagen)
	if err != nil {
		if err == sql.ErrNoRows {
			return imagen, nil // No se encuentra la imagen, devuelve un objeto vacío
		}
		log.Println("Error al obtener imagen:", err)
		return imagen, err // Retorna un error si falla la consulta
	}
	return imagen, nil // Devuelve la imagen completa si no hay error
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
func EliminarImagen(db *sql.DB, ruta string) error {
	query := `DELETE FROM producto_imagenes WHERE ruta_imagen = ?`
	_, err := db.Exec(query, ruta)
	if err != nil {
		log.Println("Error al eliminar imagen:", err)
		return err
	}
	return nil
}
