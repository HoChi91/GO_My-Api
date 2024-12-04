package controllers

import (
	database "MYAPI/config"
	"MYAPI/helper"
	"MYAPI/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var userCarts = make(map[int][]models.CartItem)

func CheckStock(bookId int, quantity int, c *gin.Context) bool {
	var stock int
	err := database.DB.QueryRow("SELECT stock FROM BOOKS WHERE id = ?", bookId).Scan(&stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recuperation du stock", "details": err.Error()})
		return false
	}
	return stock >= quantity
}

// Achat sans ajout au panier
func PurchaseBook(c *gin.Context) {
	userId, _ := c.Get("userID") // Récupérer l'ID de l'utilisateur depuis le token (authentification)

	var purchaseRequest models.PurchaseRequest
	// var purchaseHistory models.PurchaseHistory

	// Récupérer les données JSON
	if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
		helper.HandleError(c, 400, "Données Invalides", err)
		return
	}

	// Vérification du stock
	if !CheckStock(purchaseRequest.BookId, purchaseRequest.Quantity, c) {
		helper.HandleError(c, 400, "Stock insuffisant", nil)
		return
	}

	// Calcul du prix total
	var bookPrice float64
	err := database.DB.QueryRow("SELECT price FROM BOOKS WHERE id = ?", purchaseRequest.BookId).Scan(&bookPrice)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération du prix", err)
		return
	}

	totalPrice := bookPrice * float64(purchaseRequest.Quantity)

	// Insérer dans l'historique des achats
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = database.DB.Exec(`
        INSERT INTO PURCHASE_HISTORY (user, book, quantity, total_price, payment_timestamp) 
        VALUES (?, ?, ?, ?, ?)`,
		userId, purchaseRequest.BookId, purchaseRequest.Quantity, totalPrice, timestamp)

	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de l'enregistrement de l'achat", err)
		return
	}

	// Mettre à jour le stock du livre
	_, err = database.DB.Exec("UPDATE BOOKS SET stock = stock - ? WHERE id = ?", purchaseRequest.Quantity, purchaseRequest.BookId)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la mise à jour du stock", err)
		return
	}

	// purchaseHistory.User = userId.(int)
	// purchaseHistory.Book = purchaseRequest.BookId
	// purchaseHistory.Quantity = purchaseRequest.Quantity
	// purchaseHistory.Total_price = totalPrice
	// purchaseHistory.Payment_timestamp = timestamp

	helper.HandleResponse(c, 200, "Achat effectué avec succès", nil)
}

func GetAllPurchaseHistory(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, user, book, quantity, total_price, payment_timestamp FROM PURCHASE_HISTORY`)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération de l'historique complet", err)
		return
	}
	defer rows.Close()

	var history []models.PurchaseHistory
	for rows.Next() {
		var item models.PurchaseHistory
		if err := rows.Scan(&item.ID, &item.User, &item.Book, &item.Quantity, &item.Total_price, &item.Payment_timestamp); err != nil {
			helper.HandleError(c, 500, "Erreur lors de la lecture de l'historique", err)
			return
		}
		history = append(history, item)
	}

	c.JSON(http.StatusOK, history)
}

func GetPurchaseHistoryByID(c *gin.Context) {
	// Récupérer l'ID depuis les paramètres de l'URL
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if role != "admin" {
		if userID != id {
			helper.HandleError(c, 403, "Accès refusé", nil)
			return
		}

	}

	// Préparer la requête SQL
	rows, err := database.DB.Query(`SELECT user, book, quantity, total_price, payment_timestamp FROM PURCHASE_HISTORY WHERE user = ?`, id)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération de l'historique pour cet utilisateur", err)
		return
	}
	defer rows.Close()

	var history []models.PurchaseHistory
	for rows.Next() {
		var item models.PurchaseHistory
		if err := rows.Scan(&item.User, &item.Book, &item.Quantity, &item.Total_price, &item.Payment_timestamp); err != nil {
			helper.HandleError(c, 500, "Erreur lors de la récupération de l'historique pour cet utilisateur", err)
			return
		}
		history = append(history, item)
	}

	// Vérifier si aucun historique n'a été trouvé
	if len(history) == 0 {
		helper.HandleError(c, 404, "Aucun historique trouvé pour cet utilisateur", nil)
		return
	}

	c.JSON(http.StatusOK, history)
}

func AddToCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string)) // Conversion string -> int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID", "details": err.Error()})
		return
	}
	var cartItem models.CartItem

	// Récupérer les données JSON
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier le stock disponible
	if !CheckStock(cartItem.BookID, cartItem.Quantity, c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant"})
		return
	}

	// Ajouter l'article au panier
	userCarts[userId] = append(userCarts[userId], cartItem)

	c.JSON(http.StatusOK, gin.H{"message": "Article ajouté au panier avec succès"})
}

func GetCart(c *gin.Context) {
	userIDStr, _ := c.Get("userID")

	userId, err := strconv.Atoi(userIDStr.(string)) // Conversion string -> int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID", "details": err.Error()})
		return
	}
	// Vérifier si le panier existe
	cart, exists := userCarts[userId]
	if !exists {
		// Si le panier est vide, retourner une réponse avec un tableau vide
		c.JSON(http.StatusOK, gin.H{"cart": []models.CartItem{}})
		return
	}

	// Retourner le contenu du panier
	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func RemoveCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string)) // Conversion string -> int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID", "details": err.Error()})
		return
	}

	delete(userCarts, userId) // Supprime le panier de l'utilisateur

	c.JSON(http.StatusOK, gin.H{"message": "Panier vidé avec succès"})
}

func FinalizeCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string)) // Conversion string -> int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID", "details": err.Error()})
		return
	}
	cart := userCarts[userId]

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le panier est vide"})
		return
	}

	var totalAmount float64
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range cart {
		// Vérifier le stock
		if !CheckStock(item.BookID, item.Quantity, c) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant pour le livre", "bookId": item.BookID})
			return
		}

		// Récupérer le prix du livre
		var bookPrice float64
		err := database.DB.QueryRow("SELECT price FROM BOOKS WHERE id = ?", item.BookID).Scan(&bookPrice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du prix", "details": err.Error()})
			return
		}

		totalPrice := bookPrice * float64(item.Quantity)
		totalAmount += totalPrice

		// Insérer dans l'historique des achats
		_, err = database.DB.Exec(`
			INSERT INTO PURCHASE_HISTORY (user, book, quantity, total_price, payment_timestamp) 
			VALUES (?, ?, ?, ?, ?)`,
			userId, item.BookID, item.Quantity, totalPrice, timestamp)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'achat", "details": err.Error()})
			return
		}

		// Mettre à jour le stock
		_, err = database.DB.Exec(`UPDATE BOOKS SET stock = stock - ? WHERE id = ?`, item.Quantity, item.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock", "details": err.Error()})
			return
		}
	}

	// Vider le panier après finalisation
	delete(userCarts, userId)

	c.JSON(http.StatusOK, gin.H{"message": "Achat finalisé avec succès", "totalAmount": totalAmount})
}

