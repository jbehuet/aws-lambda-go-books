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
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToMarshalItem     = "failed to marshal item"

	ErrorCouldNotFetchItem  = "could not fetch item"
	ErrorCouldNotFetchItems = "could not fetch records"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotPutItem    = "could not put item"
	ErrorCouldNotUpdateItem = "could not update item"

	ErrorInvalidBookData   = "invalid book data"
	ErrorBookDoesNotExists = "Book doesn't exist"
)

type Book struct {
	UUID   string `json:"uuid"` // field name in dynamodb !case sensitive
	Title  string `json:"title"`
	Author string `json:"author"`
	Editor string `json:"editor"`
}

func FetchBook(uuid, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Book, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotFetchItem)

	}

	item := new(Book)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchBooks(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Book, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotFetchItems)
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
		return nil, errors.New(ErrorFailedToMarshalItem)
	}

	// Put item in dynamodb
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem)
	}
	return &b, nil
}

func UpdateBook(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Book, error) {
	uuid := req.QueryStringParameters["uuid"]

	var b Book
	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(ErrorInvalidBookData)
	}
	b.UUID = uuid

	currentBook, _ := FetchBook(uuid, tableName, dynaClient)
	if currentBook != nil && len(currentBook.UUID) == 0 {
		return nil, errors.New(ErrorBookDoesNotExists)
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {
				S: aws.String(uuid),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(b.Title),
			},
			":author": {
				S: aws.String(b.Author),
			},
			":editor": {
				S: aws.String(b.Editor),
			},
		},
		TableName:        aws.String(tableName),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set title = :title, author = :author, editor = :editor"),
	}

	_, err := dynaClient.UpdateItem(input)
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
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
