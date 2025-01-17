package models

// Historique des achats
type PurchaseHistory struct {
	ID                int     `json:"id"`
	User              int     `json:"user"`
	Book              int     `json:"book"`
	Quantity          int     `json:"quantity"`
	Total_price       float64 `json:"total_price"`
	Payment_timestamp string  `json:"payment_timestamp"`
}

type PurchaseHistoryWithBookDetails struct {
	ID               int     `json:"id"`
	BookID           int     `json:"book"`
	Quantity         int     `json:"quantity"`
	TotalPrice       float64 `json:"total_price"`
	PaymentTimestamp string  `json:"payment_timestamp"`
	BookTitle        string  `json:"book_title"`
	BookPrice        float64 `json:"book_price"`
}

// Représente un article dans le panier
type CartItem struct {
	BookID   int `json:"bookId"`
	Quantity int `json:"quantity"`
}

// Requête pour effectuer un achat
type PurchaseRequest struct {
	BookId   int `json:"bookId" validate:"required,min=1"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// CartItemWithDetails est un modèle de panier avec les informations détaillées sur les livres.
type CartItemWithDetails struct {
	BookID   int  `json:"bookId"`
	Quantity int  `json:"quantity"`
	Book     Book `json:"book"` // Ajout de l'objet Book
}
