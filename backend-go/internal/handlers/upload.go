package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	clients "backend-go/internal/clients"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
)

func getUserIdFromContext(c *gin.Context) int32 {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return -1
	}
	userID := userIDVal.(int32)
	return userID
}

func UploadHandler(s3Client *s3.Client, sqsClient *sqs.Client, bucketName string, queueName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIdFromContext(c)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer src.Close()

		key := fmt.Sprintf("users/%d/%s", userID, file.Filename)

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &key,
			Body:   src,
		})

		if err != nil {
			log.Printf("upload failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
			return
		}

		event := map[string]string{
			"bucket": bucketName,
			"key":    key,
			"userId": string(userID),
		}
		body, _ := json.Marshal(event)

		err = clients.SendMessage(sqsClient, queueName, string(body))
		if err != nil {
			log.Printf("failed to send SQS message: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "file uploaded successfully",
			"file":    file.Filename,
		})
	}
}

func ListFilesHandler(s3Client *s3.Client, bucketName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIdFromContext(c)
		prefix := fmt.Sprintf("users/%d/", userID)

		out, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: &bucketName,
			Prefix: &prefix,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list files"})
			return
		}

		files := []string{}
		for _, obj := range out.Contents {
			files = append(files, *obj.Key)
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	}
}
