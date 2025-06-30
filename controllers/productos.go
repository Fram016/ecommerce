package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
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
	// Obtenemos la imagen principal para cada producto
	for i, p := range productos {
		imagen, err := models.ImagenPrincipal(db, p.ID)
		if err != nil {
			http.Error(w, "Error al obtener imagen principal de productos", http.StatusInternalServerError)
			return
		}

		productos[i].ImagenPrincipal = imagen.RutaImagen // Asignamos la imagen al producto
	}

	// Obtener la sesión del usuario
	session, _ := store.Get(r, "session")
	rol, ok := session.Values["rol"].(string)
	if !ok {
		rol = ""
	}
	viewData := ViewData{
		Productos:  productos,
		UsuarioRol: rol,
	}

	if id, ok := session.Values["usuario_id"].(int); ok {
		viewData.UsuarioID = id
	}
	if correo, ok := session.Values["correo"].(string); ok {
		viewData.Correo = correo
	}
	if nombre, ok := session.Values["nombre"].(string); ok {
		viewData.Nombre = nombre
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
		err = tmpl.ExecuteTemplate(w, "productos_admin", viewData)
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
		err = tmpl.ExecuteTemplate(w, "productos_cliente", viewData)
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
	// Obtener las imágenes del producto
	imagenes, err := models.ListarImagenes(db, id)
	if err != nil {
		http.Error(w, "Error al obtener imágenes", http.StatusInternalServerError)
		return
	}

	//obtener la imagen principal del producto
	imagenPrincipal, err := models.ImagenPrincipal(db, id)
	if err != nil {
		http.Error(w, "Error al obtener la imagen principal", http.StatusInternalServerError)
		return
	}

	producto.Imagenes = imagenes
	producto.ImagenPrincipal = imagenPrincipal.RutaImagen

	viewData := ViewData{
		Producto:   producto,
		UsuarioRol: rol,
	}

	if id, ok := session.Values["usuario_id"].(int); ok {
		viewData.UsuarioID = id
	}
	if correo, ok := session.Values["correo"].(string); ok {
		viewData.Correo = correo
	}
	if nombre, ok := session.Values["nombre"].(string); ok {
		viewData.Nombre = nombre
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
		err = tmpl.ExecuteTemplate(w, "detalle_producto_admin", viewData)
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
		err = tmpl.ExecuteTemplate(w, "detalle_producto_cliente", viewData)
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
		// Obtener las imágenes del producto
		imagenes, err := models.ListarImagenes(db, id)
		if err != nil {
			http.Error(w, "Error al obtener imágenes", http.StatusInternalServerError)
			return
		}

		producto.Imagenes = imagenes

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

func SubirArchivo(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Verificar si la solicitud es POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el formulario para obtener el archivo
	err := r.ParseMultipartForm(10 << 20) // Limitar el tamaño del archivo a 10 MB
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	// Obtener el archivo del formulario
	file, handle, err := r.FormFile("imagen")
	if err != nil {
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	// obtenemos y transformamos el ID del producto
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

	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}
	ruta := "public/img/products/" + "product-" + strconv.Itoa(id) + "-" + handle.Filename
	err = os.WriteFile(ruta, data, 0666)
	if err != nil {
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}

	// Actualizar la ruta de la imagen en la base de datos
	err = models.CrearImagen(db, models.ProductoImagen{
		ProductoID: id,
		RutaImagen: ruta,
		TipoImagen: "principal",
	})

	if err != nil {
		// Verificamos si hubo un error al crear la imagen del producto
		http.Error(w, "Error al crear la imagen del producto", http.StatusInternalServerError)
		return
	}
	returnedURL := fmt.Sprintf("/admin/producto/editar?id=%d", id)
	http.Redirect(w, r, returnedURL, http.StatusSeeOther)
}

func EliminarImagen(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener la ruta de la imagen desde la URL
	ruta := r.URL.Query().Get("ruta")
	if ruta == "" {
		http.Error(w, "Ruta de imagen no proporcionada", http.StatusBadRequest)
		return
	}

	// Eliminar el archivo físicamente
	err := os.Remove(ruta)
	if err != nil {
		log.Println("Error al eliminar el archivo de la imagen:", err)
		// No hacemos return aquí porque igual intentamos eliminar de BD
	}

	// Eliminar de la base de datos
	err = models.EliminarImagen(db, ruta)
	if err != nil {
		http.Error(w, "Error al eliminar la imagen de la base de datos", http.StatusInternalServerError)
		return
	}

	// Redirigir al origen
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
