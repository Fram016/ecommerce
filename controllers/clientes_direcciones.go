package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// ListarDirecciones lista todas las direcciones de un cliente o todas las direcciones si es un admin
func ListarDirecciones(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener la sesión del usuario logueado
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error al obtener la sesión", http.StatusInternalServerError)
		return
	}

	// Verificar si el usuario está autenticado
	usuarioID := session.Values["usuario_id"]
	rol := session.Values["rol"]

	// Si no está autenticado o el ID del usuario no está en la sesión
	if usuarioID == nil {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	// Obtener el ID de usuario desde los parámetros de la URL
	usuarioIDStr := r.URL.Query().Get("id")
	var usuarioIDParam int

	// Si se proporciona un ID en la URL, lo convertimos a int
	if usuarioIDStr != "" {
		usuarioIDParam, err = strconv.Atoi(usuarioIDStr)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}
	}

	// Si el usuario es un cliente, solo puede ver sus propias direcciones
	if rol != "admin" && usuarioID.(int) != usuarioIDParam {
		http.Error(w, "No tienes permiso para acceder a las direcciones de otro usuario", http.StatusForbidden)
		return
	}

	// Si el rol es admin, mostramos todas las direcciones (o las de un usuario específico si se pasa un id en la URL)
	var direcciones []models.ClienteDireccion

	if rol == "admin" {
		// Si es admin y se proporciona un id, mostramos las direcciones de ese usuario
		if usuarioIDStr != "" {
			direcciones, err = models.ListarDirecciones(db, usuarioIDParam)
		} else {
			// Si no se proporciona un id, mostramos todas las direcciones
			direcciones, err = models.ListarDirecciones(db, 0)
		}
	} else {
		// Si no es admin, mostramos solo las direcciones del usuario logueado
		usuarioIDInt := usuarioID.(int)
		direcciones, err = models.ListarDirecciones(db, usuarioIDInt)
	}

	// Verificamos si ocurrió un error al obtener las direcciones
	if err != nil {
		http.Error(w, "Error al obtener direcciones", http.StatusInternalServerError)
		return
	}

	// Cargar la plantilla (vista) donde mostraremos las direcciones
	tmpl, err := template.ParseFiles("views/direcciones.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de direcciones", http.StatusInternalServerError)
		return
	}

	// Pasamos las direcciones a la vista
	tmpl.Execute(w, direcciones)
}

// CrearDireccion maneja la creación de una nueva dirección
func CrearDireccion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Si es una solicitud POST, obtenemos los datos del formulario
		usuarioIDStr := r.FormValue("usuario_id")
		direccion := r.FormValue("direccion")
		ciudad := r.FormValue("ciudad")
		provincia := r.FormValue("provincia")
		codigoPostal := r.FormValue("codigo_postal")
		pais := r.FormValue("pais")
		// Si el checkbox está marcado, esPrincipal será true, si no está marcado será false
		esPrincipal := false
		if r.FormValue("es_principal") == "on" {
			esPrincipal = true
		}

		// Convertir usuarioID de string a int
		usuarioID, err := strconv.Atoi(usuarioIDStr)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}

		// Insertamos la nueva dirección en la base de datos
		err = models.CrearDireccion(db, models.ClienteDireccion{
			UsuarioID:    usuarioID,
			Direccion:    direccion,
			Ciudad:       ciudad,
			Provincia:    provincia,
			CodigoPostal: codigoPostal,
			Pais:         pais,
			EsPrincipal:  esPrincipal,
		})

		if err != nil {
			http.Error(w, "Error al crear la dirección", http.StatusInternalServerError)
			return
		}

		// Redirigir a la lista de direcciones después de crear una nueva
		http.Redirect(w, r, fmt.Sprintf("/direcciones?id=%d", usuarioID), http.StatusSeeOther)
		return
	}

	// Si es GET, mostramos el formulario para crear una nueva dirección
	tmpl, err := template.ParseFiles("views/crear_direccion.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de creación", http.StatusInternalServerError)
		return
	}

	// Renderizamos el formulario
	tmpl.Execute(w, nil)
}

// ObtenerDireccion maneja la visualización de una dirección específica
func ObtenerDireccion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la dirección desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de dirección no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de dirección inválido", http.StatusBadRequest)
		return
	}

	// Llamar al modelo para obtener la dirección por ID
	direccion, err := models.ObtenerDireccion(db, id)
	if err != nil {
		http.Error(w, "Dirección no encontrada", http.StatusNotFound)
		return
	}

	// Cargar la plantilla (vista) donde mostraremos la dirección
	tmpl, err := template.ParseFiles("views/ver_direccion.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de dirección", http.StatusInternalServerError)
		return
	}

	// Pasamos la dirección a la vista
	tmpl.Execute(w, direccion)
}

// ModificarDireccion maneja la modificación de una dirección
func ModificarDireccion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la dirección desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de dirección no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de dirección inválido", http.StatusBadRequest)
		return
	}

	// Si es una solicitud GET, mostramos el formulario con los datos actuales de la dirección
	if r.Method == http.MethodGet {
		// Obtener la dirección desde la base de datos
		direccion, err := models.ObtenerDireccion(db, id)
		if err != nil {
			http.Error(w, "Dirección no encontrada", http.StatusNotFound)
			return
		}

		// Mostrar el formulario con los datos actuales de la dirección
		tmpl, err := template.ParseFiles("views/editar_direccion.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista de edición", http.StatusInternalServerError)
			return
		}

		// Pasamos la dirección a la vista para mostrar los datos actuales
		tmpl.Execute(w, direccion)
		return
	}

	// Si es una solicitud POST, actualizamos la dirección con los nuevos datos
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		usuarioIDStr := r.FormValue("usuario_id")
		direccion := r.FormValue("direccion")
		ciudad := r.FormValue("ciudad")
		provincia := r.FormValue("provincia")
		codigoPostal := r.FormValue("codigo_postal")
		pais := r.FormValue("pais")

		// Obtenemos el valor del checkbox (puede ser "" si no está marcado)
		esPrincipalStr := r.FormValue("es_principal")

		// Si el checkbox está marcado, esPrincipalStr será "on", si no está marcado será ""
		esPrincipal := false
		if esPrincipalStr == "on" {
			esPrincipal = true
		}
		// Convertir usuarioID de string a int
		usuarioID, err := strconv.Atoi(usuarioIDStr)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}

		// Crear un nuevo objeto dirección con los datos actualizados
		direccionObj := models.ClienteDireccion{
			ID:           id,
			UsuarioID:    usuarioID,
			Direccion:    direccion,
			Ciudad:       ciudad,
			Provincia:    provincia,
			CodigoPostal: codigoPostal,
			Pais:         pais,
			EsPrincipal:  esPrincipal,
		}

		// Llamar al modelo para actualizar la dirección en la base de datos
		err = models.ModificarDireccion(db, direccionObj)
		if err != nil {
			http.Error(w, "Error al modificar la dirección", http.StatusInternalServerError)
			return
		}

		// Redirigir a la página de lista de direcciones después de la modificación
		http.Redirect(w, r, fmt.Sprintf("/direcciones?id=%d", usuarioID), http.StatusSeeOther)
		return
	}
}

// EliminarDireccion maneja la eliminación de una dirección por su ID
func EliminarDireccion(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la dirección desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de dirección no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de dirección inválido", http.StatusBadRequest)
		return
	}

	// Llamar al modelo para eliminar la dirección de la base de datos
	err = models.EliminarDireccion(db, id)
	if err != nil {
		http.Error(w, "Error al eliminar la dirección", http.StatusInternalServerError)
		return
	}

	// Redirigir a la lista de direcciones después de eliminar una
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
