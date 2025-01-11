package models

type PurchaseHistory struct {
	ID                int
	User              int
	Book              int
	Quantity          int
	Total_price       float64
	Payment_timestamp string
}

type CartItem struct {
	BookID   int `json:"BookId"`
	Quantity int `json:"Quantity"`
}

type PurchaseRequest struct {
	BookId   int `json:"bookId" validate:"required,min=1"`
	Quantity int `json:"quantity" validate:"required,min=1"`
}

