package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "backend-go/internal/db"
	sqlc "backend-go/internal/db/sqlc"
	handlers "backend-go/internal/handlers"
)

var s3Client *s3.Client
var bucketName string = "file-overview-system-bucket"

func initS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-central-1"), // required by AWS SDK
	)
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Use path-style addressing
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})
	return s3Client
}

func createBucket(client *s3.Client, bucketName string) error {
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraintEuCentral1,
		},
	})
	if err != nil {
		log.Printf("failed to create bucket: %v", err)
		return err
	}
	log.Printf("Bucket %s created successfully\n", bucketName)
	return nil
}

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

	s3Client = initS3Client()
	_ = createBucket(s3Client, bucketName)

	r := gin.Default()

	r.GET("/health", handlers.HealthHandler)

	r.POST("/register", handlers.RegisterHandler(queries))

	r.POST("/login", handlers.LoginHandler(queries))

	r.POST("/upload", handlers.UploadHandler(s3Client, bucketName))

	r.POST("/files", handlers.ListFilesHandler(s3Client, bucketName))

	// Port iz .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
	fmt.Printf("Server started on http://localhost:%s\n", port)
}
