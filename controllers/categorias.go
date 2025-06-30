package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// ListarCategorias lista todas las categorías
func ListarCategorias(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Llamamos al modelo para obtener la lista de categorías desde la base de datos
	categorias, err := models.ListarCategorias(db)
	if err != nil {
		// Si ocurre un error, mostramos una respuesta 500
		http.Error(w, "Error al obtener categorías", http.StatusInternalServerError)
		return
	}

	// Cargamos la plantilla (vista) donde mostraremos las categorías
	tmpl, err := template.ParseFiles("views/admin/categorias.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de categorías", http.StatusInternalServerError)
		return
	}

	// Pasamos las categorías a la vista
	tmpl.Execute(w, categorias)
}

// CrearCategoria maneja la creación de una nueva categoría
func CrearCategoria(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Si es una solicitud POST, obtenemos los datos del formulario
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")

		// Insertamos la nueva categoría en la base de datos
		err := models.CrearCategoria(db, models.Categoria{
			Nombre:      nombre,
			Descripcion: descripcion,
		})

		if err != nil {
			http.Error(w, "Error al crear la categoría", http.StatusInternalServerError)
			return
		}

		// Redirigir a la lista de categorías después de crear una nueva
		http.Redirect(w, r, "/categorias", http.StatusSeeOther)
		return
	}

	// Si es GET, mostramos el formulario para crear una nueva categoría
	tmpl, err := template.ParseFiles("views/crear_categoria.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de creación", http.StatusInternalServerError)
		return
	}

	// Renderizamos el formulario
	tmpl.Execute(w, nil)
}

// ObtenerCategoria obtiene una categoría por su ID
func ObtenerCategoria(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la categoría desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de categoría no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Obtener la categoría desde la base de datos utilizando el ID
	categoria, err := models.ObtenerCategoria(db, id)
	if err != nil {
		http.Error(w, "Categoría no encontrada", http.StatusNotFound)
		return
	}

	// Renderizamos la vista de la categoría
	tmpl, err := template.ParseFiles("views/detalle_categoria.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de categoría", http.StatusInternalServerError)
		return
	}

	// Pasamos la categoría a la vista para que se renderice
	tmpl.Execute(w, categoria)
}

// ModificarCategoria maneja la modificación de una categoría
func ModificarCategoria(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la categoría desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de categoría no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Si es una solicitud GET, mostramos el formulario con los datos actuales de la categoría
	if r.Method == http.MethodGet {
		// Obtener la categoría desde la base de datos
		categoria, err := models.ObtenerCategoria(db, id)
		if err != nil {
			http.Error(w, "Categoría no encontrada", http.StatusNotFound)
			return
		}

		// Mostrar el formulario con los datos actuales de la categoría
		tmpl, err := template.ParseFiles("views/editar_categoria.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista de edición", http.StatusInternalServerError)
			return
		}

		// Pasamos la categoría a la vista para mostrar los datos actuales
		tmpl.Execute(w, categoria)
		return
	}

	// Si es una solicitud POST, actualizamos la categoría con los nuevos datos
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")

		// Crear un nuevo objeto categoría con los datos actualizados
		categoria := models.Categoria{
			ID:          id,
			Nombre:      nombre,
			Descripcion: descripcion,
		}

		// Llamar al modelo para actualizar la categoría en la base de datos
		err = models.ModificarCategoria(db, categoria)
		if err != nil {
			http.Error(w, "Error al modificar la categoría", http.StatusInternalServerError)
			return
		}

		// Redirigir a la página de detalles de la categoría después de la modificación
		http.Redirect(w, r, fmt.Sprintf("/categoria?id=%d", id), http.StatusSeeOther)
		return
	}
}

// EliminarCategoria maneja la eliminación de una categoría por su ID
func EliminarCategoria(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID de la categoría desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de categoría no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	// Llamar al modelo para eliminar la categoría de la base de datos
	err = models.EliminarCategoria(db, id)
	if err != nil {
		http.Error(w, "Error al eliminar la categoría", http.StatusInternalServerError)
		return
	}

	// Redirigir a la lista de categorías después de eliminar una
	http.Redirect(w, r, "/categorias", http.StatusSeeOther)
}
