package main

import (
	"ECOMMERCE-GO/db"
	"log"
	"net/http"
)
func main() {
	// vamos a conectar con la base de datos
	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// aqui podemos hacer consultas a la base de datos
	log.Println("Base de datos conectada correctamente")
	// por ejemplo, podemos hacer una consulta simple
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Versión de la base de datos:", version)
	// configuramos el manejador de rutas
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}