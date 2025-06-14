package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Categoria representa la estructura de la tabla 'categorias'
type Categoria struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// ListarCategorias devuelve todas las categorías en la base de datos
func ListarCategorias(db *sql.DB) ([]Categoria, error) {
	rows, err := db.Query(`SELECT id, nombre, descripcion FROM categorias`)
	if err != nil {
		log.Println("Error al obtener categorías:", err)
		return nil, err
	}
	defer rows.Close()

	var categorias []Categoria
	for rows.Next() {
		var categoria Categoria
		if err := rows.Scan(&categoria.ID, &categoria.Nombre, &categoria.Descripcion); err != nil {
			log.Println("Error al escanear categoría:", err)
			return nil, err
		}
		categorias = append(categorias, categoria)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return categorias, nil
}

// CrearCategoria inserta una nueva categoría en la base de datos
func CrearCategoria(db *sql.DB, categoria Categoria) error {
	query := `INSERT INTO categorias (nombre, descripcion) VALUES (?, ?)`
	_, err := db.Exec(query, categoria.Nombre, categoria.Descripcion)
	if err != nil {
		log.Println("Error al crear categoría:", err)
		return err
	}
	return nil
}

// ObtenerCategoria por ID
func ObtenerCategoria(db *sql.DB, id int) (Categoria, error) {
	var categoria Categoria
	query := `SELECT id, nombre, descripcion FROM categorias WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&categoria.ID, &categoria.Nombre, &categoria.Descripcion)
	if err != nil {
		if err == sql.ErrNoRows {
			return categoria, fmt.Errorf("categoría no encontrada")
		}
		return categoria, err
	}
	return categoria, nil
}

// ModificarCategoria actualiza la información de una categoría en la base de datos
func ModificarCategoria(db *sql.DB, categoria Categoria) error {
	query := `UPDATE categorias SET nombre = ?, descripcion = ? WHERE id = ?`
	_, err := db.Exec(query, categoria.Nombre, categoria.Descripcion, categoria.ID)
	if err != nil {
		log.Println("Error al modificar categoría:", err)
		return err
	}
	return nil
}

// EliminarCategoria elimina una categoría de la base de datos
func EliminarCategoria(db *sql.DB, id int) error {
	query := `DELETE FROM categorias WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar categoría:", err)
		return err
	}
	return nil
}
