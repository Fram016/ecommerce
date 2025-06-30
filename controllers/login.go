package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type LoginPageData struct {
	Error string
}

type ViewData struct {
	Productos  []models.Producto
	UsuarioID  int
	UsuarioRol string
	Correo     string
	Nombre     string
	Producto   models.Producto
	Categorias []models.Categoria
	Categoria  models.Categoria
	Usuarios   []models.Usuario
	Usuario   models.Usuario
}

// Crear el store de sesiones con una clave secreta
var store = sessions.NewCookieStore([]byte("VoladorDeClaveSecreta*"))

// MostrarLogin muestra el formulario de login
func MostrarLogin(w http.ResponseWriter, r *http.Request) {
	// Verificar si ya hay una sesión activa
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error al obtener la sesión", http.StatusInternalServerError)
		return
	}

	// Si ya está logueado, redirigir al dashboard
	if session.Values["usuario_id"] != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Obtener mensajes flash
	flashes := session.Flashes("error")
	var mensajeError string
	if len(flashes) > 0 {
		mensajeError = flashes[0].(string)
	}
	session.Save(r, w)

	// Mostrar el formulario de login
	tmpl, err := template.ParseFiles("views/auth/login.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de login", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, LoginPageData{Error: mensajeError})
}

// ProcesarLogin maneja el inicio de sesión de los usuarios y crea una sesión
func ProcesarLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Obtener el correo y la contraseña del formulario
		correo := r.FormValue("correo")
		claveSegura := r.FormValue("clave_segura")

		// Obtener el usuario desde la base de datos por correo
		usuario, err := models.ObtenerUsuarioPorCorreo(db, correo)
		if err != nil {
			// Si la contraseña no coincide, mostramos error
			session, _ := store.Get(r, "session")
			session.AddFlash("Usuario incorrecto", "error")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			tmpl, _ := template.ParseFiles("views/auth/login.html")
			tmpl.Execute(w, LoginPageData{Error: "Usuario incorrecto"})
			return
		}

		// Verificar que la contraseña ingresada coincida con el hash almacenado
		err = bcrypt.CompareHashAndPassword([]byte(usuario.ClaveSegura), []byte(claveSegura))
		if err != nil {
			// Si la contraseña no coincide, mostramos error
			session, _ := store.Get(r, "session")
			session.AddFlash("Contraseña incorrecta", "error")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			tmpl, _ := template.ParseFiles("views/auth/login.html")
			tmpl.Execute(w, LoginPageData{Error: "Contraseña incorrecta"})
			return
		}

		// Iniciar sesión creando una nueva sesión
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, "Error al obtener la sesión", http.StatusInternalServerError)
			return
		}

		// Guardar los datos de la sesión
		session.Values["usuario_id"] = usuario.ID
		session.Values["rol"] = usuario.Rol
		session.Values["correo"] = usuario.Correo
		session.Values["nombre"] = usuario.Nombre

		// Opcionalmente, podemos establecer opciones de expiración de la cookie
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   21600,                // La sesión expirará en 6 horas
			SameSite: http.SameSiteLaxMode, // Evitar CSRF
		}

		// Guardar la sesión
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Error al guardar la sesión", http.StatusInternalServerError)
			return
		}

		// Redirigir al dashboard o una página protegida
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Si no es POST, mostramos el formulario de login
	tmpl, err := template.ParseFiles("views/login.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de login", http.StatusInternalServerError)
		return
	}

	// Renderizamos el formulario de login
	tmpl.Execute(w, nil)
}

// CerrarSesion maneja el cierre de sesión
func CerrarSesion(w http.ResponseWriter, r *http.Request) {
	// Obtener la sesión
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error al obtener la sesión", http.StatusInternalServerError)
		return
	}

	// Eliminar la sesión
	session.Options.MaxAge = -1 // Esto elimina la sesión
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Error al guardar los cambios de la sesión", http.StatusInternalServerError)
		return
	}

	// Redirigir al login después de cerrar sesión
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// MostrarIndex maneja la solicitud para la página de inicio y carga los productos
func MostrarIndex(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener los productos desde la base de datos
	productos, err := models.ListarProductos(db)
	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}
	for i, p := range productos {
		imagen, err := models.ImagenPrincipal(db, p.ID)
		if err != nil {
			http.Error(w, "Error al obtener imagen principal del producto", http.StatusInternalServerError)
			return
		}
		productos[i].ImagenPrincipal = imagen.RutaImagen
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

	var tmpl *template.Template

	// Si el rol es admin, mostramos la vista para administradores
	if rol == "admin" {
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_admin.html",
			"views/admin/index_admin.html",
			"views/partials/footer.html",
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "index_admin", viewData)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de admin", http.StatusInternalServerError)
			return
		}

	} else {
		// Si el rol es cliente o cualquier otro, mostramos la vista para clientes
		tmpl = template.Must(template.ParseFiles(
			"views/partials/header_cliente.html", // Encabezado para cliente
			"views/index.html",                   // Cuerpo con los productos para clientes
			"views/partials/footer.html",         // Pie de página común
		))
		// Pasamos los productos a la vista y renderizamos todas las plantillas (header, cuerpo, footer)
		err = tmpl.ExecuteTemplate(w, "index", viewData)
		if err != nil {
			// Si ocurre un error al renderizar la plantilla, mostramos una respuesta 500
			http.Error(w, "Error al renderizar la vista de productos", http.StatusInternalServerError)
			return
		}
	}

}
