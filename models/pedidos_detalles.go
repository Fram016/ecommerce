package models

import (
	"database/sql"
	"log"
)

// PedidoDetalle representa la estructura de la tabla 'pedidos_detalles'
type PedidoDetalle struct {
	ID             int     `json:"id"`
	PedidoID       int     `json:"pedido_id"`
	ProductoID     int     `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Observacion    string  `json:"observacion,omitempty"`
}

// CrearDetallePedido inserta un nuevo detalle de pedido en la base de datos
func CrearDetallePedido(db *sql.DB, detalle PedidoDetalle) error {
	query := `INSERT INTO pedidos_detalles (pedido_id, producto_id, cantidad, precio_unitario, observacion) VALUES (?, ?, ?, ?, ?)`
	if detalle.Observacion != "" {
		_, err := db.Exec(query, detalle.PedidoID, detalle.ProductoID, detalle.Cantidad, detalle.PrecioUnitario, detalle.Observacion)
		if err != nil {
			log.Println("Error al crear detalle de pedido:", err)
			return err
		}
	} else {
		_, err := db.Exec(query, detalle.PedidoID, detalle.ProductoID, detalle.Cantidad, detalle.PrecioUnitario, nil)
		if err != nil {
			log.Println("Error al crear detalle de pedido:", err)
			return err
		}
	}
	return nil
}

// ListarDetallesPorPedido obtiene los detalles de un pedido específico
func ListarDetallesPorPedido(db *sql.DB, pedidoID int) ([]PedidoDetalle, error) {
	rows, err := db.Query(`SELECT id, pedido_id, producto_id, cantidad, precio_unitario, observacion FROM pedidos_detalles WHERE pedido_id = ?`, pedidoID)
	if err != nil {
		log.Println("Error al obtener detalles de pedido:", err)
		return nil, err
	}
	defer rows.Close()

	var detalles []PedidoDetalle
	for rows.Next() {
		var detalle PedidoDetalle
		if err := rows.Scan(&detalle.ID, &detalle.PedidoID, &detalle.ProductoID, &detalle.Cantidad, &detalle.PrecioUnitario, &detalle.Observacion); err != nil {
			log.Println("Error al escanear detalle de pedido:", err)
			return nil, err
		}
		detalles = append(detalles, detalle)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return detalles, nil
}

// EliminarDetallesPorPedido elimina todos los detalles de un pedido específico
func EliminarDetallesPorPedido(db *sql.DB, pedidoID int) error {
	query := `DELETE FROM pedidos_detalles WHERE pedido_id = ?`
	_, err := db.Exec(query, pedidoID)
	if err != nil {
		log.Println("Error al eliminar detalles de pedido:", err)
		return err
	}
	return nil
}
