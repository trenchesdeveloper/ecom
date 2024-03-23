package product

import (
	"database/sql"
	"github.com/trenchesdeveloper/go-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetProduct retrieves a product by its ID
func (s *Store) GetProduct(id int) (*types.Product, error) {
	row := s.db.QueryRow("SELECT * FROM products WHERE id = $1", id)
	p := &types.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.Quantity, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetProducts retrieves all products
func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		p := types.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Image, &p.Quantity, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// CreateProduct creates a new product
func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, price, description, image, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Price, product.Description, product.Image, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}
