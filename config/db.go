package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(DataSourceName string) {
	var err error

	DB, err = sql.Open("mysql", DataSourceName) // Connexion à la base de données
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données : %v", err)
	}

	err = DB.Ping() // Vérification de la connexion
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}
	log.Println("Connexion à la base de données réussie")
}
