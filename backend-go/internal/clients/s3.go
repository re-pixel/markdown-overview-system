package clients

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func InitS3Client() *s3.Client {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "eu-central-1"
	}
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("LOCALSTACK_ENDPOINT")
	}
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(endpoint)
	})
	return s3Client
}

func CreateBucket(client *s3.Client, bucketName string) error {
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
