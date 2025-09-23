package handlers

import (
	sqlc "backend-go/internal/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		user, err := queries.GetUserByEmail(c, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		sessionToken := uuid.NewString()
		expiresAt := time.Now().Add(24 * time.Hour)

		_, err = queries.CreateSession(c, sqlc.CreateSessionParams{
			UserID:       user.ID,
			SessionToken: sessionToken,
			ExpiresAt: pgtype.Timestamp{
				Time:  expiresAt,
				Valid: true},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create session token"})
		}

		c.SetCookie(
			"session_id",
			sessionToken,
			3600*24,
			"/",
			"",
			false,
			true,
		)

		c.JSON(http.StatusOK, user)
	}
}
