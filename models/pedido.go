package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Pedido representa la estructura de la tabla 'pedidos'
type Pedido struct {
	ID             int     `json:"id"`
	UsuarioID      int     `json:"usuario_id"`
	FechaPedido    string  `json:"fecha_pedido"`
	Estado         string  `json:"estado"`
	Total          float64 `json:"total"`
	DireccionEnvio int     `json:"direccion_envio"`
	Observaciones  string  `json:"observacion,omitempty"`
}

// ListarPedidos devuelve todos los pedidos en la base de datos
func ListarPedidos(db *sql.DB) ([]Pedido, error) {
	rows, err := db.Query(`SELECT id, usuario_id, fecha_pedido, estado, total, direccion_envio, observacion FROM pedidos`)
	if err != nil {
		log.Println("Error al obtener pedidos:", err)
		return nil, err
	}
	defer rows.Close()

	var pedidos []Pedido
	for rows.Next() {
		var pedido Pedido
		if err := rows.Scan(&pedido.ID, &pedido.UsuarioID, &pedido.FechaPedido, &pedido.Estado, &pedido.Total, &pedido.DireccionEnvio, &pedido.Observaciones); err != nil {
			log.Println("Error al escanear pedido:", err)
			return nil, err
		}
		pedidos = append(pedidos, pedido)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return pedidos, nil
}

// ObtenerPedido por ID
func ObtenerPedido(db *sql.DB, id int) (Pedido, error) {
	var pedido Pedido
	query := `SELECT id, usuario_id, fecha_pedido, estado, total, direccion_envio, observacion FROM pedidos WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&pedido.ID, &pedido.UsuarioID, &pedido.FechaPedido, &pedido.Estado, &pedido.Total, &pedido.DireccionEnvio, &pedido.Observaciones)
	if err != nil {
		if err == sql.ErrNoRows {
			return pedido, fmt.Errorf("pedido no encontrado")
		}
		return pedido, err
	}
	return pedido, nil
}

// CrearPedido inserta un nuevo pedido en la base de datos
func CrearPedido(db *sql.DB, pedido Pedido) error {
	query := `INSERT INTO pedidos (usuario_id, fecha_pedido, estado, total, direccion_envio, observacion) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, pedido.UsuarioID, pedido.FechaPedido, pedido.Estado, pedido.Total, pedido.DireccionEnvio, pedido.Observaciones)
	if err != nil {
		log.Println("Error al crear pedido:", err)
		return err
	}
	return nil
}

// ModificarPedido actualiza la información de un pedido en la base de datos
func ModificarPedido(db *sql.DB, pedido Pedido) error {
	query := `UPDATE pedidos SET usuario_id = ?, fecha_pedido = ?, estado = ?, total = ?, direccion_envio = ?, observacion = ? WHERE id = ?`
	_, err := db.Exec(query, pedido.UsuarioID, pedido.FechaPedido, pedido.Estado, pedido.Total, pedido.DireccionEnvio, pedido.Observaciones, pedido.ID)
	if err != nil {
		log.Println("Error al modificar pedido:", err)
		return err
	}
	return nil
}

// EliminarPedido elimina un pedido de la base de datos
func EliminarPedido(db *sql.DB, id int) error {
	query := `DELETE FROM pedidos WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar pedido:", err)
		return err
	}
	return nil
}
