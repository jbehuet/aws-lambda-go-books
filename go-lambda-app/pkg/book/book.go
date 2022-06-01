package book

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ErrorFailedToFetchRecords = "failed to fetch records"

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Editor string `json:"editor"`
}

func FetchBooks(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Book, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecords)
	}
	items := new([]Book)
	_ = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	return items, nil
}
