package main

import (
	"context"
	"jbehuet/aws-lambda-go-books/pkg/services/books"
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	dynaClient *dynamodb.Client
	s3Client   *s3.Client
)

func main() {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return
	}
	dynaClient = dynamodb.NewFromConfig(cfg)
	s3Client = s3.NewFromConfig(cfg)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		uuid := req.PathParameters["uuid"]
		if uuid != "" {
			return books.GetBook(uuid, dynaClient, s3Client)
		}
		return books.GetBooks(dynaClient)
	case "POST":
		isCoverPath, _ := regexp.Match(`.*\/cover`, []byte(req.Path))
		if isCoverPath {
			return books.UploadCover(req, s3Client)
		}
		return books.CreateBook(req, dynaClient)
	case "PUT":
		return books.UpdateBook(req, dynaClient)
	case "DELETE":
		return books.DeleteBook(req, dynaClient, s3Client)
	default:
		return utils.UnhandledMethod()
	}
}
