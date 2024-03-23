package types

type OrderStore interface {
	CreateOrder(order Order) (int, error)
	CreateOrderItem(orderItem OrderItem) error
}

type Order struct {
	ID        int     `json:"id"`
	UserID    int     `json:"userId"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	Address   string  `json:"address"`
	CreatedAt string  `json:"created_at"`
}

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

type CheckoutInput struct {
	Items []CartItem `json:"items" validate:"required"`
}
