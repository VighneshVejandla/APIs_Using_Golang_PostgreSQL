package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"zocket/cache"
	"zocket/database"
	"zocket/messaging"
	"zocket/api/handlers"
)

// RegisterRoutes defines the API routes for the application
func RegisterRoutes(r *mux.Router, db *sql.DB, redisClient *redis.Client, rabbitMQConn *amqp.Connection) {
	// Initialize messaging
	messaging.InitializeMessaging(rabbitMQConn)

	// Product routes
	r.HandleFunc("/products", GetProducts(db)).Methods("GET")
	r.HandleFunc("/product", CreateProduct(db)).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", handlers.GetProductByID(db)).Methods("GET")

	// User routes
	r.HandleFunc("/users", handlers.GetUsersHandler(db)).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserByID(db)).Methods("GET")

	// Publish message to RabbitMQ
	r.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var msg struct {
			Exchange string `json:"exchange"`
			Message  string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := messaging.PublishMessage(msg.Exchange, msg.Message); err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	// Redis cache handling
	r.HandleFunc("/cache/{key}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		key := params["key"]

		val, err := cache.GetCache(redisClient, key)
		if err != nil {
			http.Error(w, "Cache miss or error", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"key": key, "value": val})
	}).Methods("GET")
}

// GetProducts fetches all products from the database
func GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

// GetProductByID fetches a product by ID
func GetProductByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		product, err := database.GetProductByID(db, id)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// CreateProduct adds a new product to the database
func CreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product database.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := database.CreateProduct(db, product); err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// GetUsers fetches all users from the database
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email FROM users")
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []database.User
		for rows.Next() {
			var u database.User
			if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
				http.Error(w, "Error scanning users", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// GetUserByID fetches a user by ID
func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := database.GetUserByID(db, id)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// CreateUser adds a new user to the database
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user database.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := database.CreateUser(db, user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}