package handlers

import (
	sqlc "backend-go/internal/db/sqlc"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "could not hash password", http.StatusInternalServerError)
			return
		}

		newUser, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			Email: req.Email,
			Pass:  string(hashedPassword),
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newUser)
	}
}
