package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	clients "backend-go/internal/clients"
	db "backend-go/internal/db"
	sqlc "backend-go/internal/db/sqlc"
	router "backend-go/internal/router"
	worker "backend-go/internal/worker"
)

var bucketName string = "file-overview-system-bucket"
var taskQueueName string = "task-queue"
var responseQueueName string = "response-queue"

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

	s3Client := clients.InitS3Client()
	clients.CreateBucket(s3Client, bucketName)

	sqsClient := clients.InitSQSClient()
	clients.CreateQueue(sqsClient, taskQueueName)
	clients.CreateQueue(sqsClient, responseQueueName)

	if err != nil {
		log.Fatalf("failed to create SQS queue: %v", err)
	}

	worker.StartResponseWorker(sqsClient, responseQueueName)

	r := router.SetupRouter(queries, s3Client, sqsClient, bucketName, taskQueueName)

	// Port iz .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
	fmt.Printf("Server started on http://localhost:%s\n", port)
}
