package worker

import (
	"backend-go/internal/events"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ResponseMessage struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Status string `json:"status"`
	UserID string `json:"userId"`
}

type SSEMessage struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

func StartResponseWorker(sqsClient *sqs.Client, s3Client *s3.Client, queueName string, bucketName string, broadcaster *events.Broadcaster) {
	go func() {
		getOut, _ := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
			QueueName: &queueName,
		})

		queueUrl := getOut.QueueUrl

		for {
			resp, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl:            queueUrl,
				MaxNumberOfMessages: 5,
				WaitTimeSeconds:     10,
			})
			if err != nil {
				log.Printf("error receiving messages: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if len(resp.Messages) == 0 {
				continue
			}

			for _, m := range resp.Messages {
				var msg ResponseMessage
				if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
					log.Printf("failed to parse message: %v", err)
					continue
				}

				s3Resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
					Bucket: aws.String(msg.Bucket),
					Key:    aws.String(msg.Key),
				})
				if err != nil {
					log.Printf("failed to fetch from s3: %v", err)
					continue
				}

				bodyBytes, _ := io.ReadAll(s3Resp.Body)
				content := string(bodyBytes)

				sseMsg := SSEMessage{
					UserID:  msg.UserID,
					Content: content,
				}

				jsonData, _ := json.Marshal(sseMsg)

				broadcaster.Publish(string(jsonData))

				_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      queueUrl,
					ReceiptHandle: m.ReceiptHandle,
				})
				if err != nil {
					log.Printf("failed to delete message: %v", err)
				}
			}
		}
	}()
}
