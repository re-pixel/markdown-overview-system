package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	clients "backend-go/internal/clients"
	db "backend-go/internal/db"
	sqlc "backend-go/internal/db/sqlc"
	"backend-go/internal/events"
	router "backend-go/internal/router"
	worker "backend-go/internal/worker"
)

var bucketName string
var taskQueueName string
var responseQueueName string

func main() {
	_ = godotenv.Load(".env")           // attempt root .env first
	_ = godotenv.Load("backend-go.env") // fallback / legacy

	bucketName = getenvDefault("S3_BUCKET_NAME", "file-overview-system-bucket")
	taskQueueName = getenvDefault("TASK_QUEUE_NAME", "task-queue")
	responseQueueName = getenvDefault("RESPONSE_QUEUE_NAME", "response-queue")

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

	broadcaster := events.NewBroadcaster()

	worker.StartResponseWorker(sqsClient, s3Client, responseQueueName, bucketName, broadcaster)

	r := router.SetupRouter(queries, s3Client, sqsClient, bucketName, taskQueueName, broadcaster)

	port := getenvDefault("PORT", "8080")
	r.Run(":" + port)
	fmt.Printf("Server started on http://localhost:%s\n", port)
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
