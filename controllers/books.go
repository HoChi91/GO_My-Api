package controllers

import (
	database "MYAPI/config"
	"MYAPI/helper"
	"MYAPI/models"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	var rows *sql.Rows
	var err error

	role, _ := c.Get("role")

	// Requête en fonction du rôle de l'utilisateur
	var query string
	if role == "admin" {
		query = "SELECT id, title, author, publication_date, summary, stock, price FROM BOOKS"
	} else if role == "user" {
		query = "SELECT id, title, author FROM BOOKS"
	}

	// Exécution de la requête
	rows, err = database.DB.Query(query)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération des livres", err)
		return
	}

	// Parcourir les résultats et les ajouter à la slice
	for rows.Next() {
		var book models.Book

		// Scanning des résultats selon la requête
		if role == "admin" {
			if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publication_Date, &book.Summary, &book.Stock, &book.Price); err != nil {
				helper.HandleError(c, 500, "Erreur lors du scan des livres", err)
				return
			}
		} else {
			if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
				helper.HandleError(c, 500, "Erreur lors du scan des livres", err)
				return
			}
		}

		// Ajout du livre uniquement si le stock est supérieur à 0
		if book.Stock != 0 || role == "admin" { // Les admins peuvent voir les livres même avec stock à 0
			books = append(books, book)
		}
	}

	// Vérification d'erreurs durant l'itération des résultats
	if err := rows.Err(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de l'itération des livres", err)
		return
	}

	helper.HandleResponse(c, 200, "Livres", books)
}

func GetBookByID(c *gin.Context) {

	id := c.Param("id")
	var book models.Book

	err := database.DB.QueryRow("SELECT* FROM BOOKS WHERE id=?", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Publication_Date, &book.Summary, &book.Stock, &book.Price)

	if err == sql.ErrNoRows {
		helper.HandleError(c, 404, "Livre non trouvé", err)
		return
	} else if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération du livre", err)
		log.Printf("Erreur de récupération du livre : %v", err)
		return
	} else if book.Stock == 0 {
		helper.HandleError(c, 404, "Livre plus en stock", err)
	}

	helper.HandleResponse(c, 200, "Livre", book)
}

func GetBooksByAuthorID(c *gin.Context) {
	authorID := c.Param("id") // Récupère l'ID de l'auteur depuis l'URL
	var books []models.Book   // Un tableau pour stocker les livres de l'auteur

	// Requête SQL pour récupérer tous les livres de cet auteur
	rows, err := database.DB.Query("SELECT * FROM BOOKS WHERE author=?", authorID)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération des livres", err)
		log.Printf("Erreur de récupération des livres : %v", err)
		return
	}
	defer rows.Close()

	// Parcours des lignes retournées par la requête
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publication_Date, &book.Summary, &book.Stock, &book.Price); err != nil {
			helper.HandleError(c, 500, "Erreur lors de l'analyse des livres", err)
			log.Printf("Erreur d'analyse des livres : %v", err)
			return
		}
		// Ajout du livre à la liste
		books = append(books, book)
	}

	// Vérifie s'il y a des livres
	if len(books) == 0 {
		helper.HandleError(c, 404, "Aucun livre trouvé pour cet auteur", nil)
		return
	}

	// Retourne les livres de l'auteur
	helper.HandleResponse(c, 200, "Livres de l'auteur", books)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		helper.HandleError(c, 400, "Données invalides", err)
		return
	}

	// Requête pour insérer un livre dans la base de données
	_, err := database.DB.Exec("INSERT INTO BOOKS (title, author, publication_date, summary, stock, price) VALUES (?, ?, ?, ?, ?, ?)",
		book.Title, book.Author, book.Publication_Date, book.Summary, book.Stock, book.Price)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Erreur lors de l'insertion du livre",
			"details": err.Error(),
		})
		return
	}

	helper.HandleResponse(c, 200, "Livre ajouté", book)
}

func ModifyBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	// Lier les données JSON à la structure
	if err := c.ShouldBindJSON(&book); err != nil {
		helper.HandleError(c, 400, "Données invalides", err)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors du démarrage de la transaction", err)
		return
	}

	// Mettre à jour les informations du livre dans la base de données
	result, err := tx.Exec("UPDATE BOOKS SET title=?, author=?, publication_date=?, summary=?, stock=?, price=? WHERE id=?", book.Title, book.Author, book.Publication_Date, book.Summary, book.Stock, book.Price, id)
	if err != nil {
		tx.Rollback()
		helper.HandleError(c, 500, "Erreur lors de la mise à jour du livre", err)
		return
	}

	// Vérifier si aucune ligne n'a été affectée
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		helper.HandleError(c, 500, "Erreur lors de la vérification des lignes affectées", err)
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		helper.HandleError(c, 404, "Auteur introuvable", err)
		return
	}

	if err := tx.Commit(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de la validation de la transaction", err)
		return
	}

	helper.HandleResponse(c, 200, "Livre mis à jour", nil)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	tx, err := database.DB.Begin()
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors du démarrage de la transaction", err)
		return
	}

	_, err = tx.Exec("DELETE FROM BOOKS WHERE id=?", id)
	if err != nil {
		tx.Rollback()
		helper.HandleError(c, 500, "Erreur lors de la suppression du livre", err)
		return
	}

	if err := tx.Commit(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de la validation de la transaction", err)
		return
	}

	helper.HandleResponse(c, 200, "Livre supprimé", nil)
}
