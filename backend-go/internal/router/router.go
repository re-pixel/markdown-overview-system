package router

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	handlers "backend-go/internal/handlers"
	middleware "backend-go/internal/middleware"

	sqlc "backend-go/internal/db/sqlc"
)

func SetupRouter(queries *sqlc.Queries, s3Client *s3.Client, sqsClient *sqs.Client, bucketName string, queueName string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // frontend URL
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", handlers.HealthHandler)

	r.POST("/register", handlers.RegisterHandler(queries))

	r.POST("/login", handlers.LoginHandler(queries))

	auth := r.Group("/")
	auth.Use(middleware.SessionMiddleware(queries))
	{
		auth.POST("/upload", handlers.UploadHandler(s3Client, sqsClient, bucketName, queueName))

		auth.POST("/files", handlers.ListFilesHandler(s3Client, bucketName))
	}

	return r
}
