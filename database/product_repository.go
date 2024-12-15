package database

import (
	"database/sql"
	"fmt"
)

// Product represents a product in the database
type Product struct {
	// ID    int
	// Name  string
	// Price float64
	ID          int
    Name        string
    Description string
    Price       float64
    CreatedAt   string
    UpdatedAt   string
}

// CreateProduct inserts a new product into the database
func CreateProduct(db *sql.DB, product Product) error {
	query := "INSERT INTO products (name, price) VALUES ($1, $2)"
	_, err := db.Exec(query, product.Name, product.Price)
	if err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}
	return nil
}

// GetProductByID fetches a product by its ID from the database
func GetProductByID(db *sql.DB, id int) (*Product, error) {
    query := "SELECT id, name, description, price, created_at, updated_at FROM products WHERE id = $1"
    row := db.QueryRow(query, id)

    var product Product
    if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt); err != nil {
        return nil, fmt.Errorf("could not find product with ID %d: %v", id, err)
    }

    return &product, nil
}