package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	// vamos a configurar el enlace con mariadb
	dns := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	// probar la conexion con la bd, haciendo ping
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Conectado de forma exitosa a la base de datos")
	return db, nil
}