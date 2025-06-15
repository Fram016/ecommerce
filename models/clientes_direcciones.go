package models

import (
	"database/sql"
	"log"
)

// ClienteDireccion representa la estructura de la tabla 'clientes_direcciones'
type ClienteDireccion struct {
	ID           int    `json:"id"`
	UsuarioID    int    `json:"usuario_id"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Provincia    string `json:"provincia"`
	CodigoPostal string `json:"codigo_postal"`
	Pais         string `json:"pais"`
	EsPrincipal  bool   `json:"es_principal"`
}

// ListarDirecciones devuelve todas las direcciones de un cliente
func ListarDirecciones(db *sql.DB, usuarioID int) ([]ClienteDireccion, error) {
	rows, err := db.Query(`SELECT id, usuario_id, direccion, ciudad, provincia, codigo_postal, pais, es_principal FROM clientes_direcciones WHERE usuario_id = ?`, usuarioID)
	if err != nil {
		log.Println("Error al obtener direcciones:", err)
		return nil, err
	}
	defer rows.Close()

	var direcciones []ClienteDireccion
	for rows.Next() {
		var direccion ClienteDireccion
		if err := rows.Scan(&direccion.ID, &direccion.UsuarioID, &direccion.Direccion, &direccion.Ciudad, &direccion.Provincia, &direccion.CodigoPostal, &direccion.Pais, &direccion.EsPrincipal); err != nil {
			log.Println("Error al escanear dirección:", err)
			return nil, err
		}
		direcciones = append(direcciones, direccion)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return direcciones, nil
}

// CrearDireccion inserta una nueva dirección en la base de datos
func CrearDireccion(db *sql.DB, direccion ClienteDireccion) error {
	query := `INSERT INTO clientes_direcciones (usuario_id, direccion, ciudad, provincia, codigo_postal, pais, es_principal) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, direccion.UsuarioID, direccion.Direccion, direccion.Ciudad, direccion.Provincia, direccion.CodigoPostal, direccion.Pais, direccion.EsPrincipal)
	if err != nil {
		log.Println("Error al crear dirección:", err)
		return err
	}
	return nil
}

// ModificarDireccion actualiza la información de una dirección
func ModificarDireccion(db *sql.DB, direccion ClienteDireccion) error {
	query := `UPDATE clientes_direcciones SET direccion = ?, ciudad = ?, provincia = ?, codigo_postal = ?, pais = ?, es_principal = ? WHERE id = ?`
	_, err := db.Exec(query, direccion.Direccion, direccion.Ciudad, direccion.Provincia, direccion.CodigoPostal, direccion.Pais, direccion.EsPrincipal, direccion.ID)
	if err != nil {
		log.Println("Error al modificar dirección:", err)
		return err
	}
	return nil
}

// EliminarDireccion elimina una dirección de la base de datos
func EliminarDireccion(db *sql.DB, id int) error {
	query := `DELETE FROM clientes_direcciones WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar dirección:", err)
		return err
	}
	return nil
}

// ObtenerDireccion obtiene una dirección de cliente por ID
func ObtenerDireccion(db *sql.DB, id int) (ClienteDireccion, error) {
	var direccion ClienteDireccion
	query := `SELECT id, usuario_id, direccion, ciudad, provincia, codigo_postal, pais, es_principal FROM clientes_direcciones WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&direccion.ID, &direccion.UsuarioID, &direccion.Direccion, &direccion.Ciudad, &direccion.Provincia, &direccion.CodigoPostal, &direccion.Pais, &direccion.EsPrincipal)
	if err != nil {
		if err == sql.ErrNoRows {
			return direccion, nil // Si no encuentra la dirección, la devuelve vacía
		}
		log.Println("Error al obtener dirección:", err)
		return direccion, err
	}
	return direccion, nil
}
