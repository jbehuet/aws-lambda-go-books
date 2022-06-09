package main

import (
	"context"
	"jbehuet/aws-lambda-go-books/pkg/services/covers"
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func main() {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return
	}
	s3Client = s3.NewFromConfig(cfg)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return utils.UnhandledMethod()
	case "POST":
		return covers.UploadCover(req, s3Client)
	default:
		return utils.UnhandledMethod()
	}
}
