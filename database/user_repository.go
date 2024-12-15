package database

import (
	"database/sql"
	"fmt"
)

// User represents a user in the database
type User struct {
	ID    int
	Username  string
	Password string
	Email string
}

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user User) error {
	// Ensure that the password is also being passed to the query
	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// GetUsers fetches all users from the database
func GetUsers(db *sql.DB) ([]User, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		// Ensure you scan the right columns as per your table schema (username, email)
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %v", err)
	}

	return users, nil
}

// GetUserByID fetches a user by their ID from the database
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := "SELECT id, username, email FROM users WHERE id = $1"
	row := db.QueryRow(query, id)

	var user User
	// Ensure you scan the correct columns
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	return &user, nil
}