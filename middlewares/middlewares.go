// middlewares/middlewares.go
package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Crear el store de sesiones con una clave secreta
var Store = sessions.NewCookieStore([]byte("VoladorDeClaveSecreta*"))

// Middleware de autenticación
func Autenticado(w http.ResponseWriter, r *http.Request) (int, *sessions.Session) {
	session, err := Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Error al obtener sesión", http.StatusInternalServerError)
		return http.StatusInternalServerError, nil
	}

	// Verificar si el usuario está autenticado
	if session.Values["usuario_id"] == nil {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return http.StatusUnauthorized, nil
	}

	return http.StatusOK, session
}

// Middleware para verificar el rol
func VerificarRol(rolRequerido string) func(http.ResponseWriter, *http.Request) bool {
	return func(w http.ResponseWriter, r *http.Request) bool {
		status, session := Autenticado(w, r)
		if status != http.StatusOK {
			return false
		}

		// Obtener el rol del usuario desde la sesión
		rol := session.Values["rol"].(string)
		if rol != rolRequerido {
			http.Error(w, "Acceso denegado", http.StatusForbidden)
			return false
		}

		return true
	}
}
