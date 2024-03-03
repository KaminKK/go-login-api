package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Register represents the JSON structure for registration
type Register struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jsonInput Register
	err := json.NewDecoder(r.Body).Decode(&jsonInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(jsonInput.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password, firstname, lastname) VALUES ($1, $2, $3, $4)", jsonInput.Username, string(hashedPassword), jsonInput.Firstname, jsonInput.Lastname)
	if err != nil {
		// Handle database errors
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Account created successfully"))
}

// Login represents the JSON structure for login
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jsonInput Login
	err := json.NewDecoder(r.Body).Decode(&jsonInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", jsonInput.Username).Scan(&hashedPassword)
	if err != nil {
		// Handle database errors or user not found
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(jsonInput.Password))
	if err != nil {
		// Incorrect password
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Password matches, generate JWT token
	token, err := GenerateJWT(jsonInput.Username)
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	// Respond with success message and JWT token
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
