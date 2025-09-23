package middleware

import (
	db "backend-go/internal/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
			return
		}

		session, err := queries.GetSession(c, cookie)
		if err != nil || session.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			return
		}
		c.Set("user_id", session.UserID)
		c.Next()
	}
}
