package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Producto representa la estructura de la tabla 'productos'
type Producto struct {
	ID            int     `json:"id"`
	Nombre        string  `json:"nombre"`
	Descripcion   string  `json:"descripcion"`
	Precio        float64 `json:"precio"`
	Stock         int     `json:"stock"`
	CategoriaID   int     `json:"categoria_id"`
	FechaAgregado string  `json:"fecha_agregado"`
	ImagenPrincipal string `json:"imagen_principal"`
	Imagenes    []ProductoImagen `json:"imagenes"`
}

// ListarProductos devuelve todos los productos en la base de datos
func ListarProductos(db *sql.DB) ([]Producto, error) {
	rows, err := db.Query(`SELECT id, nombre, descripcion, precio, stock, categoria_id, fecha_agregado FROM productos`)
	if err != nil {
		log.Println("Error al obtener productos:", err)
		return nil, err
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var producto Producto
		if err := rows.Scan(&producto.ID, &producto.Nombre, &producto.Descripcion, &producto.Precio, &producto.Stock, &producto.CategoriaID, &producto.FechaAgregado); err != nil {
			log.Println("Error al escanear producto:", err)
			return nil, err
		}
		productos = append(productos, producto)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return productos, nil
}

// ObtenerProducto por ID
func ObtenerProducto(db *sql.DB, id int) (Producto, error) {
	var producto Producto
	query := `SELECT id, nombre, descripcion, precio, stock, categoria_id, fecha_agregado FROM productos WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&producto.ID, &producto.Nombre, &producto.Descripcion, &producto.Precio, &producto.Stock, &producto.CategoriaID, &producto.FechaAgregado)
	if err != nil {
		if err == sql.ErrNoRows {
			return producto, fmt.Errorf("producto no encontrado")
		}
		return producto, err
	}
	return producto, nil
}

// CrearProducto inserta un nuevo producto en la base de datos
func CrearProducto(db *sql.DB, producto Producto) error {
	query := `INSERT INTO productos (nombre, descripcion, precio, stock, categoria_id) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, producto.Nombre, producto.Descripcion, producto.Precio, producto.Stock, producto.CategoriaID)
	if err != nil {
		log.Println("Error al crear producto:", err)
		return err
	}
	return nil
}

// ModificarProducto actualiza la información de un producto en la base de datos
func ModificarProducto(db *sql.DB, producto Producto) error {
	query := `UPDATE productos SET nombre = ?, descripcion = ?, precio = ?, stock = ?, categoria_id = ? WHERE id = ?`
	_, err := db.Exec(query, producto.Nombre, producto.Descripcion, producto.Precio, producto.Stock, producto.CategoriaID, producto.ID)
	if err != nil {
		log.Println("Error al modificar producto:", err)
		return err
	}
	return nil
}

// EliminarProducto elimina un producto de la base de datos
func EliminarProducto(db *sql.DB, id int) error {
	query := `DELETE FROM productos WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar producto:", err)
		return err
	}
	return nil
}
