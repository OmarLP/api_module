package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OmarLP/api_module/pkg/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	// cargar las variables de entorno desde el archivo .env
	_ = godotenv.Load()

	// inicializamos el logger para mensajes de log
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	// puedo mostrar un mensaje de conección exitosa con el nombre de la base de datos
	fmt.Printf("Conected to darabase: %v - %s\n", db, db.Migrator().CurrentDatabase())

	port := ":8080"
	address := "" + port
	fmt.Printf("Server is running on http://localhost%s\n", address)

	log.Fatal(http.ListenAndServe(address, nil))

}
