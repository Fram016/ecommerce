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
	viewData := ViewData{
		Categorias: categorias,
	}
	tmpl := template.Must(template.ParseFiles(
		"views/partials/header_admin.html",
		"views/admin/categorias.html",
		"views/partials/footer.html",
	))
	// Pasamos las categorías a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
	err = tmpl.ExecuteTemplate(w, "categorias", viewData)
	if err != nil {
		// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
		http.Error(w, "Error al renderizar la vista de categorías", http.StatusInternalServerError)
		return
	}

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
		http.Redirect(w, r, "/admin/categorias", http.StatusSeeOther)
		return
	}

	viewData := ViewData{}
	tmpl := template.Must(template.ParseFiles(
		"views/partials/header_admin.html",
		"views/admin/crear_categoria.html",
		"views/partials/footer.html",
	))
	// Pasamos las categorías a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
	err := tmpl.ExecuteTemplate(w, "crear_categoria", viewData)
	if err != nil {
		// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
		http.Error(w, "Error al renderizar la vista de categorías", http.StatusInternalServerError)
		return
	}
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

	viewData := ViewData{
		Categoria: categoria,
	}
	tmpl := template.Must(template.ParseFiles(
		"views/partials/header_admin.html",
		"views/admin/detalle_categoria.html",
		"views/partials/footer.html",
	))
	// Pasamos la categoría a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
	err = tmpl.ExecuteTemplate(w, "ver_categoria", viewData)
	if err != nil {
		// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
		http.Error(w, "Error al renderizar la vista de categorías", http.StatusInternalServerError)
		return
	}
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

		viewData := ViewData{
			Categoria: categoria,
		}
		tmpl := template.Must(template.ParseFiles(
			"views/partials/header_admin.html",
			"views/admin/editar_categoria.html",
			"views/partials/footer.html",
		))
		// Pasamos la categoría a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "editar_categoria", viewData)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de categorías", http.StatusInternalServerError)
			return
		}
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
			http.Redirect(w, r, fmt.Sprintf("/admin/categoria?id=%d", id), http.StatusSeeOther)
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
	http.Redirect(w, r, "/admin/categorias", http.StatusSeeOther)
}
