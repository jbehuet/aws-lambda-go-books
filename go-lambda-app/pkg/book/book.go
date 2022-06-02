package book

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecords  = "failed to fetch records"
	ErrorInvalidBookData       = "invalid book data"
	ErrorCouldNotMarshalItem   = "could not marshal item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item error"
)

type Book struct {
	Title  string `json:"Title"` // field name in dynamodb !case sensitive
	Author string `json:"Author"`
	Editor string `json:"Editor"`
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
	_ = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	return items, nil
}

func CreateBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Book, error) {

	// Create a book from request body
	var b Book
	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(ErrorInvalidBookData)
	}

	// Create dynamodb attributeValues
	av, err := dynamodbattribute.MarshalMap(b)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	// Put item in dynamodb
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	fmt.Println(input)

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &b, nil
}
