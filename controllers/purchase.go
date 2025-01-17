package controllers

import (
	database "MYAPI/config"
	"MYAPI/helper"
	"MYAPI/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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

	if err := validator.New().Struct(purchaseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides", "details": err.Error()})
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
	id, _ := c.Get("userID")

	// Requête pour récupérer l'historique des achats avec les détails des livres
	rows, err := database.DB.Query(`
		SELECT p.id, p.book, p.quantity, p.total_price, p.payment_timestamp, b.title, b.price
		FROM PURCHASE_HISTORY p
		JOIN BOOKS b ON p.book = b.id
		WHERE p.user = ?`, id)
	if err != nil {
		helper.HandleError(c, 500, "Erreur lors de la récupération de l'historique pour cet utilisateur", err)
		return
	}
	defer rows.Close()

	var history []models.PurchaseHistoryWithBookDetails
	for rows.Next() {
		var item models.PurchaseHistoryWithBookDetails
		if err := rows.Scan(&item.ID, &item.BookID, &item.Quantity, &item.TotalPrice, &item.PaymentTimestamp, &item.BookTitle, &item.BookPrice); err != nil {
			helper.HandleError(c, 500, "Erreur lors de la lecture de l'historique", err)
			return
		}
		history = append(history, item)
	}

	if len(history) == 0 {
		helper.HandleError(c, 404, "Aucun historique trouvé pour cet utilisateur", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}


func AddToCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string)) // Conversion string -> int
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID utilisateur invalide", "details": err.Error()})
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

	// Vérifier si le panier de l'utilisateur existe déjà
	if _, exists := userCarts[userId]; !exists {
		userCarts[userId] = []models.CartItem{}
	}

	// Vérifier si le livre est déjà dans le panier
	var updatedCart []models.CartItem
	itemFound := false
	for _, item := range userCarts[userId] {
		if item.BookID == cartItem.BookID {
			// Si le livre est déjà dans le panier, on met à jour la quantité
			item.Quantity += cartItem.Quantity
			updatedCart = append(updatedCart, item)
			itemFound = true
		} else {
			updatedCart = append(updatedCart, item)
		}
	}

	// Si l'article n'était pas dans le panier, on l'ajoute
	if !itemFound {
		updatedCart = append(updatedCart, cartItem)
	}

	// Mettre à jour le panier de l'utilisateur
	userCarts[userId] = updatedCart

	c.JSON(http.StatusOK, gin.H{"message": "Article ajouté/quantité mise à jour avec succès"})
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
		c.JSON(http.StatusOK, gin.H{"cart": []models.CartItemWithDetails{}})
		return
	}

	// Enrichir le panier avec les informations sur les livres
	var cartWithBooks []models.CartItemWithDetails

	for _, item := range cart {
		// Requête pour récupérer les détails du livre (ex: titre et prix)
		var book models.Book
		err := database.DB.QueryRow("SELECT id, title, price FROM BOOKS WHERE id = ?", item.BookID).Scan(&book.ID, &book.Title, &book.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des informations du livre", "details": err.Error()})
			return
		}

		// Ajouter les détails du livre au panier
		cartItemWithDetails := models.CartItemWithDetails{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Book:     book,
		}

		// Ajouter l'élément enrichi dans le panier
		cartWithBooks = append(cartWithBooks, cartItemWithDetails)
	}

	// Retourner le panier avec les livres enrichis
	c.JSON(http.StatusOK, gin.H{"cart": cartWithBooks})
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

// FinalizeCart
func FinalizeCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userId, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID", "details": err.Error()})
		return
	}
	cart := userCarts[userId]

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le panier est vide"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la transaction", "details": err.Error()})
		return
	}

	var totalAmount float64
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range cart {
		// Vérification du stock
		var stock int
		err := tx.QueryRow("SELECT stock FROM BOOKS WHERE id = ?", item.BookID).Scan(&stock)
		if err != nil || stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant pour le livre", "bookId": item.BookID})
			return
		}

		// Récupération du prix
		var bookPrice float64
		err = tx.QueryRow("SELECT price FROM BOOKS WHERE id = ?", item.BookID).Scan(&bookPrice)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du prix", "details": err.Error()})
			return
		}

		totalPrice := bookPrice * float64(item.Quantity)
		totalAmount += totalPrice

		// Insérer dans l'historique
		_, err = tx.Exec(
			"INSERT INTO PURCHASE_HISTORY (user, book, quantity, total_price, payment_timestamp) VALUES (?, ?, ?, ?, ?)",
			userId, item.BookID, item.Quantity, totalPrice, timestamp,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'enregistrement de l'achat", "details": err.Error()})
			return
		}

		// Mettre à jour le stock
		_, err = tx.Exec("UPDATE BOOKS SET stock = stock - ? WHERE id = ?", item.Quantity, item.BookID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock", "details": err.Error()})
			return
		}
	}

	// Valider la transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la validation de la transaction", "details": err.Error()})
		return
	}

	// Vider le panier
	delete(userCarts, userId)

	c.JSON(http.StatusOK, gin.H{"message": "Achat finalisé avec succès", "totalAmount": totalAmount})
}
