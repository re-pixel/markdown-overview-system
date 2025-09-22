package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "backend-go/internal/db"
	sqlc "backend-go/internal/db/sqlc"
	handlers "backend-go/internal/handlers"
)

func main() {
	err := godotenv.Load("backend-go.env")
	if err != nil {
		log.Println(".env file not found")
	}

	conn, err := db.Connect()

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	defer conn.Close()
	fmt.Println("Connected to NeonDB")

	queries := sqlc.New(conn)

	r := gin.Default()

	r.GET("/health", handlers.HealthHandler)

	r.POST("/register", handlers.RegisterHandler(queries))

	r.POST("/login", handlers.LoginHandler(queries))

	// Port iz .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
	fmt.Printf("Server started on http://localhost:%s\n", port)
}
