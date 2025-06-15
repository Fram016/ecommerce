package controllers

import (
	"ECOMMERCE-GO/models"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// ListarUsuarios lista todos los usuarios
func ListarUsuarios(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Llamamos al modelo para obtener la lista de usuarios desde la base de datos
	usuarios, err := models.ListarUsuarios(db)
	if err != nil {
		// Si ocurre un error, mostramos una respuesta 500
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	// Cargamos la plantilla (vista) donde mostraremos los usuarios
	tmpl, err := template.ParseFiles("views/usuarios.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista de usuarios", http.StatusInternalServerError)
		return
	}

	// Pasamos los usuarios a la vista
	tmpl.Execute(w, usuarios)
}

// CrearUsuario maneja la creación de un nuevo usuario con encriptación de la contraseña usando bcrypt
func CrearUsuario(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Si es una solicitud POST, obtenemos los datos del formulario
		Correo := r.FormValue("correo")
		Nombre := r.FormValue("nombre")
		ClaveSegura := r.FormValue("clave_segura")
		Rol := r.FormValue("rol")

		// Encriptamos la contraseña usando bcrypt
		claveHash, err := bcrypt.GenerateFromPassword([]byte(ClaveSegura), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
			return
		}

		// Insertamos el nuevo usuario con la contraseña encriptada en la base de datos
		err = models.CrearUsuario(db, models.Usuario{
			Correo:      Correo,
			Nombre:      Nombre,
			ClaveSegura: string(claveHash), // Guardamos la contraseña encriptada
			Rol:         Rol,
		})

		if err != nil {
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			return
		}

		// Redirigir a la lista de usuarios después de crear uno nuevo
		http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
		return
	}

	// Si es GET, mostramos el formulario para crear un nuevo usuario
	tmpl, err := template.ParseFiles("views/crear_usuario.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de creación", http.StatusInternalServerError)
		return
	}

	// Renderizamos el formulario
	tmpl.Execute(w, nil)
}

// RegistrarUsuario maneja el registro de un nuevo usuario con rol "cliente" y encriptación de contraseña usando bcrypt
func RegistrarUsuario(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Si el método es POST, procesamos el formulario
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		Correo := r.FormValue("correo")
		Nombre := r.FormValue("nombre")
		ClaveSegura := r.FormValue("clave_segura")

		// Asignamos el rol siempre como "cliente"
		Rol := "cliente"

		// Encriptamos la contraseña usando bcrypt
		claveHash, err := bcrypt.GenerateFromPassword([]byte(ClaveSegura), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error al encriptar la contraseña", http.StatusInternalServerError)
			return
		}

		// Insertamos el nuevo usuario con rol "cliente" y contraseña encriptada en la base de datos
		err = models.CrearUsuario(db, models.Usuario{
			Correo:      Correo,
			Nombre:      Nombre,
			ClaveSegura: string(claveHash), // Guardamos la contraseña encriptada
			Rol:         Rol,
		})

		if err != nil {
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			return
		}

		// Redirigir al login después de registrar el usuario
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Si es GET, mostramos el formulario para crear un nuevo usuario cliente
	tmpl, err := template.ParseFiles("views/registrarse.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario de creación", http.StatusInternalServerError)
		return
	}

	// Renderizamos el formulario
	tmpl.Execute(w, nil)
}

// obtener un usuario por su ID
func ObtenerUsuario(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del usuario desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Obtener el usuario desde la base de datos
	usuario, err := models.ObtenerUsuario(db, id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Renderizar la vista del usuario
	tmpl, err := template.ParseFiles("views/detalle_usuario.html")
	if err != nil {
		http.Error(w, "Error al cargar la vista del usuario", http.StatusInternalServerError)
		return
	}

	// Pasamos el usuario a la vista para que se renderice
	tmpl.Execute(w, usuario)
}

// ModificarUsuario maneja la modificación de un usuario
func ModificarUsuario(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del usuario desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Si es una solicitud GET, mostramos el formulario con los datos actuales del usuario
	if r.Method == http.MethodGet {
		// Obtener el usuario desde la base de datos
		usuario, err := models.ObtenerUsuario(db, id)
		if err != nil {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
			return
		}

		// Mostrar el formulario con los datos actuales del usuario
		tmpl, err := template.ParseFiles("views/editar_usuario.html")
		if err != nil {
			http.Error(w, "Error al cargar la vista de edición", http.StatusInternalServerError)
			return
		}

		// Pasamos el usuario a la vista para mostrar los datos actuales
		tmpl.Execute(w, usuario)
		return
	}

	// Si es una solicitud POST, actualizamos el usuario con los nuevos datos
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		correo := r.FormValue("correo")
		nombre := r.FormValue("nombre")
		claveSegura := r.FormValue("clave_segura")
		rol := r.FormValue("rol")

		// Crear un nuevo objeto usuario con los datos actualizados
		usuario := models.Usuario{
			ID:          id,
			Correo:      correo,
			Nombre:      nombre,
			ClaveSegura: claveSegura,
			Rol:         rol,
		}

		// Llamar al modelo para actualizar el usuario en la base de datos
		err = models.ModificarUsuario(db, usuario)
		if err != nil {
			http.Error(w, "Error al modificar el usuario", http.StatusInternalServerError)
			return
		}

		// Redirigir a la página de detalles del usuario después de la modificación
		http.Redirect(w, r, fmt.Sprintf("/usuario?id=%d", id), http.StatusSeeOther)
		return
	}
}

// EliminarUsuario maneja la eliminación de un usuario por su ID
func EliminarUsuario(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Obtener el ID del usuario desde los parámetros de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	// Convertir el ID de string a int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Llamar al modelo para eliminar el usuario de la base de datos
	err = models.EliminarUsuario(db, id)
	if err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	// Redirigir a la lista de usuarios después de eliminar uno
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
