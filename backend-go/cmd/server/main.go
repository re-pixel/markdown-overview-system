package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	db "backend-go/internal/db"
)

func main() {
	err := godotenv.Load("../../backend-go.env")
	if err != nil {
		log.Println(".env file not found")
	}

	conn, err := db.Connect()

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to NeonDB")

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running!"))
	}).Methods("GET")

	// Port iz .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
