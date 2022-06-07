package book

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

var (
	ErrorFailedToFetchRecords  = "failed to fetch records"
	ErrorInvalidBookData       = "invalid book data"
	ErrorCouldNotMarshalItem   = "could not marshal item"
	ErrorCouldNotDynamoPutItem = "could not dynamo put item error"
	ErrorCouldNotDeleteItem    = "could not delete item"
)

type Book struct {
	UUID   string `json:"uuid"` // field name in dynamodb !case sensitive
	Title  string `json:"Title"`
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
	b.UUID = uuid.NewString()

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

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &b, nil
}

func DeleteBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	uuid := req.QueryStringParameters["uuid"]
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
