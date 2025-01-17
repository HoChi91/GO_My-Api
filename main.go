package main

import (
	database "MYAPI/config"
	router "MYAPI/routes"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// DataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", // Variable d'environnement pour connexion MariaDB
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_NAME"),
	// )

	DataSourceName := "API:API@tcp(localhost:3306)/library"

	port := os.Getenv("PORT") //Variable d'environnement pour le serveur
	if port == "" {
		port = "4000" //port par defaut
	}

	database.InitDB(DataSourceName) //Initialisation de la base de données

	// Router
	Router := router.SetupRoutes()

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: Router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Echec lors du démarrage du serveur")
	} else {
		log.Printf("Démarrage du serveur sur le port %s", port)
	}

}
