package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func FetchSummaryHandler(s3Client *s3.Client, bucketName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userName := c.Query("userName")
		fileName := c.Query("file")

		if userName == "" || fileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing username or file"})
			return
		}

		summaryKey := userName + "/" + fileName + "_overview.txt"

		out, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    &summaryKey,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
			return
		}
		defer out.Body.Close()

		body, err := io.ReadAll(out.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read summary"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"file":    summaryKey,
			"summary": string(body),
		})
	}
}
