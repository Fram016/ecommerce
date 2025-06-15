package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Usuario representa la estructura de la tabla 'usuarios'
type Usuario struct {
	ID            int    `json:"id"`
	Correo        string `json:"correo"`
	Nombre        string `json:"nombre"`
	ClaveSegura   string `json:"clave_segura"`
	Rol           string `json:"rol"`
	FechaRegistro string `json:"fecha_registro"`
}

// ListarUsuarios devuelve todos los usuarios en la base de datos
func ListarUsuarios(db *sql.DB) ([]Usuario, error) {
	rows, err := db.Query(`SELECT id, correo, nombres, clave_segura, rol, fecha_registro FROM usuarios`)
	if err != nil {
		log.Println("Error al obtener usuarios:", err)
		return nil, err
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var usuario Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Correo, &usuario.Nombre, &usuario.ClaveSegura, &usuario.Rol, &usuario.FechaRegistro); err != nil {
			log.Println("Error al escanear usuario:", err)
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error al iterar sobre filas:", err)
		return nil, err
	}
	return usuarios, nil
}

// CrearUsuario inserta un nuevo usuario en la base de datos
func CrearUsuario(db *sql.DB, usuario Usuario) error {
	query := `INSERT INTO usuarios (correo, nombres, clave_segura, rol, fecha_registro) VALUES (?, ?, ?, ?, NOW())`
	_, err := db.Exec(query, usuario.Correo, usuario.Nombre, usuario.ClaveSegura, usuario.Rol)
	if err != nil {
		log.Println("Error al crear usuario:", err)
		return err
	}
	return nil
}

// ObtenerUsuarioPorCorreo devuelve un usuario basado en su correo
func ObtenerUsuarioPorCorreo(db *sql.DB, correo string) (Usuario, error) {
	var usuario Usuario
	query := `SELECT id, correo, nombres, clave_segura, rol, fecha_registro FROM usuarios WHERE correo = ?`
	err := db.QueryRow(query, correo).Scan(&usuario.ID, &usuario.Correo, &usuario.Nombre, &usuario.ClaveSegura, &usuario.Rol, &usuario.FechaRegistro)
	if err != nil {
		if err == sql.ErrNoRows {
			return usuario, fmt.Errorf("usuario no encontrado")
		}
		return usuario, err
	}
	return usuario, nil
}

func ObtenerUsuario(db *sql.DB, id int) (Usuario, error) {
	var usuario Usuario
	query := `SELECT id, correo, nombres, clave_segura, rol, fecha_registro FROM usuarios WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&usuario.ID, &usuario.Correo, &usuario.Nombre, &usuario.ClaveSegura, &usuario.Rol, &usuario.FechaRegistro)
	if err != nil {
		if err == sql.ErrNoRows {
			return usuario, fmt.Errorf("usuario no encontrado")
		}
		return usuario, err
	}
	return usuario, nil
}

// ModificarUsuario actualiza la información de un usuario en la base de datos
func ModificarUsuario(db *sql.DB, usuario Usuario) error {
	query := `UPDATE usuarios SET correo = ?, nombres = ?, clave_segura = ?, rol = ? WHERE id = ?`
	_, err := db.Exec(query, usuario.Correo, usuario.Nombre, usuario.ClaveSegura, usuario.Rol, usuario.ID)
	if err != nil {
		log.Println("Error al modificar usuario:", err)
		return err
	}
	return nil
}

// EliminarUsuario elimina un usuario de la base de datos
func EliminarUsuario(db *sql.DB, id int) error {
	query := `DELETE FROM usuarios WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Println("Error al eliminar usuario:", err)
		return err
	}
	return nil
}
