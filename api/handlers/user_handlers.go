package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"database/sql"
	"zocket/database" // Only import database to access models
	"github.com/gorilla/mux" // Import mux for routing
	"strconv"
)

// GetUsersHandler fetches all users from the database
func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch users from the database
		users, err := database.GetUsers(db) // Ensure `database` is only used for database access
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
			return
		}

		// Set response header and send the users list as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}


// CreateUser handles the request to create a new user
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user database.User
		// Decode the JSON request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Create user in database
		if err := database.CreateUser(db, user); err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}

		// Send success response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// GetUsers handles the request to fetch users from the database
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch users from the database
		users, err := database.GetUsers(db)
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		// Send users data in response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// GetUserByID handles the request to fetch a user by ID from the database
func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from URL
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Fetch user from the database
		user, err := database.GetUserByID(db, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch user: %v", err), http.StatusInternalServerError)
			return
		}

		// Send user data in response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}