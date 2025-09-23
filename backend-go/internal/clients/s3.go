package clients

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func InitS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-central-1"), // required by AWS SDK
	)
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Use path-style addressing
		o.BaseEndpoint = aws.String("http://localhost:4566")
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
