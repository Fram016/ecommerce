package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

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

	// Mostrar el formulario de login
	tmpl, err := template.ParseFiles("views/auth/login.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de login", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
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
			// Si el usuario no existe, mostramos error
			http.Error(w, "Usuario incorrectos", http.StatusUnauthorized)
			return
		}

		// Verificar que la contraseña ingresada coincida con el hash almacenado
		err = bcrypt.CompareHashAndPassword([]byte(usuario.ClaveSegura), []byte(claveSegura))
		if err != nil {
			// Si la contraseña no coincide, mostramos error
			http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
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

		// Opcionalmente, podemos establecer opciones de expiración de la cookie
		session.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 3600, // La sesión expirará en 1 hora
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

	// Cargar la plantilla de la página de inicio
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de inicio", http.StatusInternalServerError)
		return
	}

	// Pasar los productos a la vista
	tmpl.Execute(w, productos)
}
