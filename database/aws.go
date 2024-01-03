package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"log"
)

var (
	sesClient *ses.Client
	s3Client  *s3.Client
	cfg       aws.Config
	err       error
)

func SetupAWSClient() error {
	cfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	s3Client = s3.NewFromConfig(cfg)

	return nil
}

func GetS3Uploader() *manager.Uploader {
	return manager.NewUploader(s3Client)
}

func GetSesClient() *ses.Client {
	if sesClient == nil {
		sesClient = ses.NewFromConfig(cfg)
	}

	return sesClient
}
