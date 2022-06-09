package main

import (
	"context"
	"jbehuet/aws-lambda-go-books/pkg/services/books"
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynaClient *dynamodb.Client

func main() {
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return
	}
	dynaClient = dynamodb.NewFromConfig(cfg)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return books.GetBookOrBooks(req, dynaClient)
	case "POST":
		return books.CreateBook(req, dynaClient)
	case "PUT":
		return books.UpdateBook(req, dynaClient)
	case "DELETE":
		return books.DeleteBook(req, dynaClient)
	default:
		return utils.UnhandledMethod()
	}
}
