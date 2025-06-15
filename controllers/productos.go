package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// MostrarProductos lista todos los productos
func MostrarProductos(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Llamamos al modelo para obtener la lista de productos desde la base de datos
	productos, err := models.ListarProductos(db)
	if err != nil {
		// Si ocurre un error, mostramos una respuesta 500
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session")
	rol, ok := session.Values["rol"].(string)
	if !ok {
		rol = ""
	}

	// Si el usuario es admin, redirigimos a la página de administración
	if rol == "admin" && !strings.HasPrefix(r.URL.Path, "/admin") {
		// Redirigimos solo si el usuario está intentando acceder a la ruta /productos o cualquier otra que no sea /admin
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
		return
	}

	// Determinamos cuál plantilla renderizar dependiendo del rol
	var tmpl *template.Template

	// Si el rol es admin, mostramos la vista para administradores
	if rol == "admin" {
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_admin.html",
			"views/admin/productos_admin.html",
			"views/partials/footer.html",
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "productos_admin", productos)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}

	} else {
		// Si el rol es cliente o cualquier otro, mostramos la vista para clientes
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_cliente.html",   // Encabezado para cliente
			"views/cliente/productos_cliente.html", // Cuerpo con los productos para clientes
			"views/partials/footer.html",           // Pie de página común
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "productos_cliente", productos)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}
	}

}

// CrearProducto maneja la creación de un nuevo producto
func CrearProducto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Si es una solicitud POST, obtenemos los datos del formulario
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")

		// Convertimos precio de string a float64
		precioStr := r.FormValue("precio")
		precio, err := strconv.ParseFloat(precioStr, 64) // Convierte de string a float64
		if err != nil {
			http.Error(w, "Error al convertir el precio", http.StatusBadRequest)
			return
		}

		// Convertimos stock de string a int
		stockStr := r.FormValue("stock")
		stock, err := strconv.Atoi(stockStr) // Convierte de string a int
		if err != nil {
			http.Error(w, "Error al convertir el stock", http.StatusBadRequest)
			return
		}

		// Convertimos categoriaID de string a int (asumiendo que categoría es un ID numérico)
		categoriaIDStr := r.FormValue("categoria_id")
		categoriaID, err := strconv.Atoi(categoriaIDStr)
		if err != nil {
			http.Error(w, "Error al convertir el ID de la categoría", http.StatusBadRequest)
			return
		}

		// Insertamos el nuevo producto en la base de datos
		err = models.CrearProducto(db, models.Producto{
			Nombre:      nombre,
			Descripcion: descripcion,
			Precio:      precio,      // Ahora es de tipo float64
			Stock:       stock,       // Ahora es de tipo int
			CategoriaID: categoriaID, // Ahora es de tipo int
		})

		if err != nil {
			http.Error(w, "Error al crear el producto", http.StatusInternalServerError)
			return
		}

		// Redirigir a la lista de productos después de crear uno nuevo
		http.Redirect(w, r, "/productos", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"views/partials/header_admin.html", // Encabezado para admin
		"views/admin/crear_producto.html",  // Cuerpo con el formulario para crear productos
		"views/partials/footer.html",       // Pie de página común
	))
	// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
	err := tmpl.ExecuteTemplate(w, "crear_producto", nil)
	if err != nil {
		// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
		http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
		return
	}
}

// ObtenerProducto obtiene un producto por su ID
func ObtenerProducto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del producto desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de producto no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}
	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session")
	rol, ok := session.Values["rol"].(string)
	if !ok {
		rol = ""
	}
	// Si el usuario es admin, redirigimos a la página de administración
	if rol == "admin" && !strings.HasPrefix(r.URL.Path, "/admin") {
		// Redirigimos solo si el usuario está intentando acceder a la ruta /productos o cualquier otra que no sea /admin
		http.Redirect(w, r, "/admin/productos", http.StatusSeeOther)
		return
	}

	// Obtener el producto desde la base de datos utilizando el ID
	producto, err := models.ObtenerProducto(db, id)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	// Determinamos cuál plantilla renderizar dependiendo del rol
	var tmpl *template.Template

	// Si el rol es admin, mostramos la vista para administradores
	if rol == "admin" {
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_admin.html",
			"views/admin/detalle_producto_admin.html",
			"views/partials/footer.html",
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "detalle_producto_admin", producto)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}

	} else {
		// Si el rol es cliente o cualquier otro, mostramos la vista para clientes
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_cliente.html",          // Encabezado para cliente
			"views/cliente/detalle_producto_cliente.html", // Cuerpo con los productos para clientes
			"views/partials/footer.html",                  // Pie de página común
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "detalle_producto_cliente", producto)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}
	}
}

// ModificarProducto maneja la modificación de un producto
func ModificarProducto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del producto desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de producto no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	// Si es una solicitud GET, mostramos el formulario con los datos actuales del producto
	if r.Method == http.MethodGet {
		// Obtener el producto desde la base de datos
		producto, err := models.ObtenerProducto(db, id)
		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}

		tmpl := template.Must(template.ParseFiles(
			"views/partials/header_admin.html", // Encabezado para admin
			"views/admin/editar_producto.html", // Cuerpo con el formulario para editar productos
			"views/partials/footer.html",       // Pie de página común
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "editar_producto", producto)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}
		return
	}

	// Si es una solicitud POST, actualizamos el producto con los nuevos datos
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		nombre := r.FormValue("nombre")
		descripcion := r.FormValue("descripcion")
		precioStr := r.FormValue("precio")
		precio, err := strconv.ParseFloat(precioStr, 64)
		if err != nil {
			http.Error(w, "Precio inválido", http.StatusBadRequest)
			return
		}

		stockStr := r.FormValue("stock")
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "Stock inválido", http.StatusBadRequest)
			return
		}

		categoriaIDStr := r.FormValue("categoria_id")
		categoriaID, err := strconv.Atoi(categoriaIDStr)
		if err != nil {
			http.Error(w, "Categoría inválida", http.StatusBadRequest)
			return
		}

		// Crear un nuevo objeto producto con los datos actualizados
		producto := models.Producto{
			ID:          id,
			Nombre:      nombre,
			Descripcion: descripcion,
			Precio:      precio,
			Stock:       stock,
			CategoriaID: categoriaID,
		}

		// Llamar al modelo para actualizar el producto en la base de datos
		err = models.ModificarProducto(db, producto)
		if err != nil {
			http.Error(w, "Error al modificar el producto", http.StatusInternalServerError)
			return
		}

		// Redirigir a la página de detalles del producto después de la modificación
		http.Redirect(w, r, fmt.Sprintf("/producto?id=%d", id), http.StatusSeeOther)
		return
	}
}

// EliminarProducto maneja la eliminación de un producto por su ID
func EliminarProducto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del producto desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de producto no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	// Llamar al modelo para eliminar el producto de la base de datos
	err = models.EliminarProducto(db, id)
	if err != nil {
		http.Error(w, "Error al eliminar el producto", http.StatusInternalServerError)
		return
	}

	// Redirigir a la lista de productos después de eliminar uno
	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}
