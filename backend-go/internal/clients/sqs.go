package clients

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func InitSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-central-1"),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})

	return client
}

func CreateQueue(client *sqs.Client, queueName string) (string, error) {
	out, err := client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	if err != nil {
		return "", err
	}
	log.Printf("Queue %s created successfully at %s\n", queueName, *out.QueueUrl)
	return *out.QueueUrl, nil
}

func SendMessage(client *sqs.Client, queueName string, messageBody string) error {
	getOut, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return err
	}

	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    getOut.QueueUrl,
		MessageBody: &messageBody,
	})

	return err
}

//func RecieveMessages(client *sqs.Client, queueUrl)
