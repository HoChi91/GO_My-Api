package controllers

import (
	database "MYAPI/config"
	"MYAPI/helper"
	"MYAPI/middleware"
	"MYAPI/models"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	// Définir une slice pour stocker les résultats
	var users []models.User

	// Récupérer les utilisateurs depuis la base de données
	rows, err := database.DB.Query("SELECT id, email FROM USERS")
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération des utilisateurs", err)
		log.Printf("Échec de la récupération des utilisateurs : %v", err)
		return
	}
	defer rows.Close()

	// Parcourir les résultats et les ajouter à la slice
	for rows.Next() {
		var user models.User
		// Scan des résultats de la base de données dans la struct user
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			helper.HandleError(c, 500, "Erreur lors du scan des utilisateurs", err)
			return
		}
		// Ajouter l'utilisateur à la slice
		users = append(users, user)
	}

	// Vérifier s'il y a eu une erreur pendant l'itération
	if err := rows.Err(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de l'itération des résultats des utilisateurs", err)
		log.Printf("Erreur pendant l'itération des résultats des utilisateurs : %v", err)
		return
	}

	// Répondre avec les utilisateurs sous forme de JSON
	helper.HandleResponse(c, 200, "Utilisateurs", users)
}

func GetUserByID(c *gin.Context) {
	// Récupérer l'ID utilisateur depuis le contexte (déjà défini par le middleware AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		helper.HandleError(c, 401, "ID utilisateur manquant dans le contexte", nil)
		return
	}

	// Vérifier que l'ID est une chaîne de caractères
	userIDStr, ok := userID.(string)
	if !ok {
		helper.HandleError(c, 401, "ID utilisateur mal formé", nil)
		return
	}

	// Récupérer l'utilisateur depuis la base de données
	var user models.User
	err := database.DB.QueryRow("SELECT id, first_name, last_name, email, role FROM USERS WHERE id=?", userIDStr).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)

	if err == sql.ErrNoRows {
		helper.HandleError(c, 404, "Utilisateur non trouvé", err)
		return
	} else if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération de l'utilisateur", err)
		log.Printf("Erreur de récupération de l'utilisateur : %v", err)
		return
	}

	// Retourner les informations de l'utilisateur
	helper.HandleResponse(c, 200, "Utilisateur", user)
}

func CreateUser(c *gin.Context) {

	// Définir un struct pour recevoir les données JSON
	var user models.User

	// Lier les données JSON à la structure
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, 400, "Données Invalides", err)
		return
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors du hashing du mot de pass", err)
		return
	}

	// Déterminer le rôle : vérifier si un utilisateur existe déjà
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM USERS").Scan(&count)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la vérification des utilisateurs existants", err)
		return
	}

	role := "user" // Rôle par défaut
	if count == 0 {
		role = "admin" // Premier utilisateur devient admin
	}

	user.Role = role

	// Requête pour insérer l'utilisateur dans la base de données
	_, err = database.DB.Exec("INSERT INTO USERS (first_name, last_name, email, password, role) VALUES (?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, string(hashedPassword), user.Role)
	if err != nil {
		helper.HandleError(c, 500, "Erreur d'insertion", err)
		return
	}

	helper.HandleResponse(c, 200, "Utilisateur créer avec succès !", nil)
}

func ModifyUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.HandleError(c, 400, "ID invalide", err)
		return
	}

	var user models.User

	// Lier les données JSON à la structure
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, 400, "Données Invalides", err)
		return
	}

	// Vérifier si l'utilisateur avec l'ID existe
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM USERS WHERE id = ?", id).Scan(&count)
	if err != nil {
		helper.HandleError(c, 500, "Erreur de base de données", err)
		return
	}

	if count == 0 {
		// Si l'utilisateur n'existe pas
		helper.HandleError(c, 404, "Utilisateur non trouvé", err)
		return
	}

	// Exécuter la requête de mise à jour
	result, err := database.DB.Exec("UPDATE USERS SET first_name=?, last_name=?, email=?, role=? WHERE id=?",
		user.FirstName, user.LastName, user.Email, strings.ToLower(user.Role), idInt)

	// Vérifier s'il y a une erreur d'exécution
	if err != nil {
		helper.HandleError(c, 500, "Données Invalides", err)
		return
	}

	// Vérifier si aucune ligne n'a été modifiée
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		helper.HandleError(c, 500, "Données Invalides", err)
		return
	}

	if rowsAffected == 0 {
		// Si aucune ligne n'est affectée, l'ID n'existe pas dans la base de données
		helper.HandleError(c, 404, "Aucune modification requise", err)
		return
	}

	helper.HandleResponse(c, 200, "Utilisateur mis a jour", user)
}

func DeleteUser(c *gin.Context) {

	id := c.Param("id")

	// Vérifier si l'utilisateur existe avant de tenter de le supprimer
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM USERS WHERE id = ?", id).Scan(&count)
	if err != nil {
		helper.HandleError(c, 500, "Erreur de base de données", err)
		return
	}

	if count == 0 {
		// Si l'utilisateur n'existe pas, retourner une erreur 404
		helper.HandleError(c, 404, "Utilisateur non trouvé", err)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors du démarrage de la transaction", err)
		return
	}

	_, err = tx.Exec("DELETE FROM USERS WHERE id=?", id)
	if err != nil {
		tx.Rollback() // Annuler la transaction en cas d'erreur
		helper.HandleError(c, 500, "Erreur lors de la suppression de l'utilisateur", err)
		return
	}

	if err := tx.Commit(); err != nil {
		helper.HandleError(c, 500, "Erreur lors de la validation de la transaction", err)
		return
	}

	helper.HandleResponse(c, 200, "Utilisateur supprimé", nil)
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		helper.HandleError(c, 400, "Données invalides", err)
		return
	}

	var storedPassword string
	var userID string
	var userRole string
	err := database.DB.QueryRow("SELECT id, password, role FROM USERS WHERE email=?", user.Email).Scan(&userID, &storedPassword, &userRole)
	if err != nil {
		helper.HandleError(c, 400, "Identifiants incorrects", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		helper.HandleError(c, 400, "Identifiants incorrects", err)
		return
	}

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID,
		"email": user.Email,
		"role":  userRole,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expire dans 24 heures
	})

	tokenString, err := Token.SignedString(middleware.JwtKey)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la génération du token", err)
		return
	}

	// Ajouter le token dans un cookie HTTP-only
	c.SetCookie("token", tokenString, 86400, "/", "", false, true) // 86400 = 24 heures

	helper.HandleResponse(c, 200, "Connexion réussi !", tokenString)
}
