// internal/worker/response_worker.go
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type ResponseMessage struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Status string `json:"status"`
	UserID string `json:"userId"`
}

func StartResponseWorker(client *sqs.Client, queueName string) {
	go func() {
		getOut, _ := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
			QueueName: &queueName,
		})

		for {
			resp, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl:            getOut.QueueUrl,
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

				fmt.Printf("Summary ready for user %s at s3://%s/%s (status: %s)\n",
					msg.UserID, msg.Bucket, msg.Key, msg.Status)

				_, err := client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      getOut.QueueUrl,
					ReceiptHandle: m.ReceiptHandle,
				})
				if err != nil {
					log.Printf("failed to delete message: %v", err)
				}
			}
		}
	}()
}
