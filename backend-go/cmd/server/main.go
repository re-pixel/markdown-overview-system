package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	db "backend-go/internal/db"
	sqlc "backend-go/internal/db/sqlc"
	handle "backend-go/internal/handlers"
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

	r := mux.NewRouter()

	r.HandleFunc("/health", handle.HealthHandler).Methods("GET")

	r.HandleFunc("/register", handle.RegisterHandler(queries)).Methods("POST")

	// Port iz .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
