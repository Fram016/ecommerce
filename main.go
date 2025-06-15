package main

import (
	"ECOMMERCE-GO/controllers" // Reemplaza con tu paquete correcto
	"ECOMMERCE-GO/db"
	"ECOMMERCE-GO/middlewares" // Reemplaza con tu paquete correcto
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Conectar a la base de datos
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer db.Close()

	// Crear el enrutador de Mux
	r := mux.NewRouter()
	///////////////////////////////// Ruta de Inicio //////////////////////////////////
	// Ruta para la página de inicio (index)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controllers.MostrarIndex(w, r, db)
	}).Methods("GET")

	///////////////////////////////// Rutas para Login //////////////////////////////////

	// Ruta para mostrar el formulario de login (GET)
	r.HandleFunc("/login", controllers.MostrarLogin).Methods("GET")

	// Ruta para procesar el login (POST)
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.ProcesarLogin(w, r, db)
	}).Methods("POST")

	// Ruta para cerrar sesión
	r.HandleFunc("/cerrar-sesion", controllers.CerrarSesion).Methods("GET")

	///////////////////////////////// Subrouter para Admin //////////////////////////////////
	// Crear un subrouter para las rutas protegidas por el rol "admin"
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if middlewares.VerificarRol("admin")(w, r) {
				next.ServeHTTP(w, r)
			}
		})
	})

	///////////////////////////////// Rutas para Productos //////////////////////////////////
	// Ruta para la lista de productos (GET)
	adminRouter.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
		controllers.MostrarProductos(w, r, db)
	}).Methods("GET")

	// Ruta para la lista de productos (GET)
	r.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
		controllers.MostrarProductos(w, r, db)
	}).Methods("GET")

	// Ruta para crear un producto (GET y POST)
	adminRouter.HandleFunc("/producto/crear", func(w http.ResponseWriter, r *http.Request) {
		controllers.CrearProducto(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para ver un producto (GET por ID)
	adminRouter.HandleFunc("/producto", func(w http.ResponseWriter, r *http.Request) {
		controllers.ObtenerProducto(w, r, db)
	}).Methods("GET")

	// Ruta para editar un producto (GET y POST)
	adminRouter.HandleFunc("/producto/editar", func(w http.ResponseWriter, r *http.Request) {
		controllers.ModificarProducto(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para eliminar un producto (GET)
	adminRouter.HandleFunc("/producto/eliminar", func(w http.ResponseWriter, r *http.Request) {
		controllers.EliminarProducto(w, r, db)
	}).Methods("GET")

	///////////////////////////////// Rutas para Categorías //////////////////////////////////

	// Ruta para la lista de categorías (GET)
	adminRouter.HandleFunc("/categorias", func(w http.ResponseWriter, r *http.Request) {
		controllers.ListarCategorias(w, r, db)
	}).Methods("GET")

	// Ruta para crear una categoría (GET y POST)
	adminRouter.HandleFunc("/categoria/crear", func(w http.ResponseWriter, r *http.Request) {
		controllers.CrearCategoria(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para ver una categoría (GET por ID)
	adminRouter.HandleFunc("/categoria", func(w http.ResponseWriter, r *http.Request) {
		controllers.ObtenerCategoria(w, r, db)
	}).Methods("GET")

	// Ruta para editar una categoría (GET y POST)
	adminRouter.HandleFunc("/categoria/editar", func(w http.ResponseWriter, r *http.Request) {
		controllers.ModificarCategoria(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para eliminar una categoría (GET)
	adminRouter.HandleFunc("/categoria/eliminar", func(w http.ResponseWriter, r *http.Request) {
		controllers.EliminarCategoria(w, r, db)
	}).Methods("GET")
	///////////////////////////////// Rutas para Usuarios //////////////////////////////////
	// Ruta para la lista de usuarios (GET)
	adminRouter.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		controllers.ListarUsuarios(w, r, db)
	}).Methods("GET")
	// Ruta para crear un usuario (GET y POST)
	adminRouter.HandleFunc("/usuario/crear", func(w http.ResponseWriter, r *http.Request) {
		controllers.CrearUsuario(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para ver un usuario (GET por ID)
	adminRouter.HandleFunc("/usuario", func(w http.ResponseWriter, r *http.Request) {
		controllers.ObtenerUsuario(w, r, db)
	}).Methods("GET")

	// Ruta para editar un usuario (GET y POST)
	adminRouter.HandleFunc("/usuario/editar", func(w http.ResponseWriter, r *http.Request) {
		controllers.ModificarUsuario(w, r, db)
	}).Methods("GET", "POST")

	// Ruta para eliminar un usuario (GET)
	adminRouter.HandleFunc("/usuario/eliminar", func(w http.ResponseWriter, r *http.Request) {
		controllers.EliminarUsuario(w, r, db)
	}).Methods("GET")

	///////////////////////////////// Ruta para Registrar Usuario (Siempre Cliente) //////////////////////////////////
	// Ruta para el formulario de registro
	r.HandleFunc("/registrarse", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegistrarUsuario(w, r, db)
	}).Methods("GET", "POST")

	///////////////////////////////// Rutas para Direcciones //////////////////////////////////
	// Ruta para la lista de direcciones (GET)
	r.HandleFunc("/direcciones", func(w http.ResponseWriter, r *http.Request) {
		controllers.ListarDirecciones(w, r, db)
	}).Methods("GET")
	// Ruta para crear una dirección (GET y POST)
	r.HandleFunc("/direccion/crear", func(w http.ResponseWriter, r *http.Request) {
		controllers.CrearDireccion(w, r, db)
	}).Methods("GET", "POST")
	// Ruta para ver una dirección (GET por ID)
	r.HandleFunc("/direccion", func(w http.ResponseWriter, r *http.Request) {
		controllers.ObtenerDireccion(w, r, db)
	}).Methods("GET")
	// Ruta para editar una dirección (GET y POST)
	r.HandleFunc("/direccion/editar", func(w http.ResponseWriter, r *http.Request) {
		controllers.ModificarDireccion(w, r, db)
	}).Methods("GET", "POST")
	// Ruta para eliminar una dirección (GET)
	r.HandleFunc("/direccion/eliminar", func(w http.ResponseWriter, r *http.Request) {
		controllers.EliminarDireccion(w, r, db)
	}).Methods("POST")

	// Servir el servidor en el puerto 8080
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
