package controllers

import (
	database "MYAPI/config"
	"MYAPI/helper"
	"MYAPI/models"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CreateAuthor(c *gin.Context) {

	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil { // Lier les données JSON à la structure

		helper.HandleError(c, 400, "Données invalides", err)
		return
	}

	_, err := database.DB.Exec("INSERT INTO AUTHORS (name, birth_date, description) VALUES (?, ?, ?)", author.Name, author.Birth_Date, author.Description) //Inserer dans la base de données
	if err != nil {

		helper.HandleError(c, 500, "Erreur lors de l'insertion de l'auteur", err)
	}

	helper.HandleResponse(c, 200, "Auteur ajouté avec succès", author) //Réponse en cas de succès
}

func GetAuthors(c *gin.Context) {
	var author models.Author

	rows, err := database.DB.Query("SELECT id, name FROM AUTHORS") // Récupérer les auteurs depuis la base de données
	if err != nil {

		helper.HandleError(c, 500, "Erreur lors de la récupération des auteurs", err)
		return
	}
	defer rows.Close()

	var authors []models.Author

	for rows.Next() { // Parcourir les résultats et les ajouter à la slice

		if err := rows.Scan(&author.ID, &author.Name); err != nil { // Scan des résultats de la base de données dans la struct user

			helper.HandleError(c, 500, "Erreur lors du scan des auteurs", err)
			return
		}

		authors = append(authors, author) // Ajouter l'utilisateur à la slice
	}

	if err := rows.Err(); err != nil { // Vérifier s'il y a eu une erreur pendant l'itération

		helper.HandleError(c, 500, "Erreur lors de la récupération des auteurs", err)
		return
	}

	helper.HandleResponse(c, 200, "Auteurs", authors) //Réponse en cas de succès
}

func GetAuthorByID(c *gin.Context) {

	id := c.Param("id")

	var author models.Author

	err := database.DB.QueryRow("SELECT id, name, birth_date, description FROM AUTHORS WHERE id=?", id).Scan(&author.ID, &author.Name, &author.Birth_Date, &author.Description) // Récupérer l'auteur depuis la base de données

	if err == sql.ErrNoRows {

		helper.HandleError(c, 404, "Auteur non trouvé", err)
		return
	} else if err != nil {

		helper.HandleError(c, 500, "Erreur lors de la récupération de l'auteur", err)
		return
	}

	helper.HandleResponse(c, 200, "auteur", author)
}

func ModifyAuthor(c *gin.Context) {
	id := c.Param("id")
	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil { // Lier les données JSON à la structure

		helper.HandleError(c, 400, "Données Invalides", err)
		return
	}

	tx, err := database.DB.Begin() //Début de la transaction
	if err != nil {

		helper.HandleError(c, 500, "Erreur lors du démarrage de la transaction", err)
		return
	}

	result, err := tx.Exec("UPDATE AUTHORS SET name=?, birth_date=?, description=? WHERE id=?", author.Name, author.Birth_Date, author.Description, id) // Mettre à jour les informations de l'auteur dans la base de données
	if err != nil {

		tx.Rollback()
		helper.HandleError(c, 400, "Erreur lors de la mise à jour de l'auteur", err)
		return
	}

	rowsAffected, err := result.RowsAffected() // Vérifier si aucune ligne n'a été affectée
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

	if err := tx.Commit(); err != nil { //Fin de la transaction

		helper.HandleError(c, 500, "Erreur lors de la validation de la transaction", err)
		return
	}

	helper.HandleResponse(c, 200, "Auteur mis à jour avec succès", nil)
}

func DeleteAuthor(c *gin.Context) {
	id := c.Param("id")

	tx, err := database.DB.Begin()
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors du démarrage de la transaction", err)
		return
	}

	_, err = tx.Exec("DELETE FROM AUTHORS WHERE id=?", id)
	if err != nil {
		tx.Rollback()
		helper.HandleError(c, 500, "Erreur lors de la suppression de l'auteur", err)
		return
	}

	if err := tx.Commit(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de la validation de la transaction", err)
		return
	}

	helper.HandleResponse(c, 200, "Auteur supprimé", nil)
}
