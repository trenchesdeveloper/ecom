package types

type ProductStore interface {
	GetProduct(id int) (*Product, error)
	GetProducts() ([]Product, error)
	CreateProduct(product Product) error
	GetProductByIDs(ids []int) ([]Product, error)
	UpdateProduct(product Product) error
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
}

type CreateProductInput struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image"`
	Quantity    int     `json:"quantity" validate:"required"`
}
