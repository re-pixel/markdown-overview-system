package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadHandler(s3Client *s3.Client, bucketName string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &file.Filename,
			Body:   src,
		})

		if err != nil {
			log.Printf("upload failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": ""})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "file uploaded successfully",
			"file":    file.Filename,
		})
	}
}

func ListFilesHandler(s3Client *s3.Client, bucketName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: &bucketName,
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
