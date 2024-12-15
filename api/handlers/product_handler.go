package handlers

import (
    "database/sql"  // Add this import
    "encoding/json"
    "net/http"
    "zocket/database"
)

// GetProducts handles the request to fetch products from the database
func GetProducts(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Fetch products from the database
        rows, err := db.Query("SELECT id, name, price FROM products")
        if err != nil {
            http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var products []database.Product
        for rows.Next() {
            var p database.Product
            if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
                http.Error(w, "Error scanning products", http.StatusInternalServerError)
                return
            }
            products = append(products, p)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(products)
    }
}


// GetProductByID handles the request to fetch a product by its ID
func GetProductByID(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get the product ID from URL parameters
        idParam := mux.Vars(r)["id"]
        id, err := strconv.Atoi(idParam)
        if err != nil {
            http.Error(w, "Invalid product ID", http.StatusBadRequest)
            return
        }

        // Fetch product from database
        product, err := database.GetProductByID(db, id)
        if err != nil {
            http.Error(w, fmt.Sprintf("Product not found: %v", err), http.StatusNotFound)
            return
        }

        // Return the product as JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(product)
    }
}