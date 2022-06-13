package books

import (
	"context"
	"encoding/json"
	"errors"
	"jbehuet/aws-lambda-go-books/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/google/uuid"
)

const tableName = "books"

func FetchByUUID(uuid string, dynaClient *dynamodb.Client) (*Book, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"UUID": &types.AttributeValueMemberS{Value: uuid},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(context.TODO(), input)
	if err != nil {
		return nil, errors.New(utils.ErrorCouldNotFetchItem)

	}

	item := new(Book)
	err = attributevalue.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(utils.ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func Fetch(dynaClient *dynamodb.Client) (*[]Book, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(context.TODO(), input)
	if err != nil {
		return nil, errors.New(utils.ErrorCouldNotFetchItems)
	}
	items := new([]Book)
	_ = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	return items, nil
}

func Create(req events.APIGatewayProxyRequest, dynaClient *dynamodb.Client) (*Book, error) {

	// Create a book from request body
	var b Book
	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(utils.ErrorInvalidBookData)
	}
	b.UUID = uuid.NewString()

	// Create dynamodb attributeValues
	av, err := attributevalue.MarshalMap(b)
	if err != nil {
		return nil, errors.New(utils.ErrorFailedToMarshalItem)
	}

	// Put item in dynamodb
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(context.TODO(), input)
	if err != nil {
		return nil, errors.New(utils.ErrorCouldNotPutItem)
	}
	return &b, nil
}

func Update(req events.APIGatewayProxyRequest, dynaClient *dynamodb.Client) (*Book, error) {
	uuid := req.PathParameters["uuid"]
	if uuid == "" {
		return nil, errors.New(utils.UUIDMissing)
	}

	var b Book
	if err := json.Unmarshal([]byte(req.Body), &b); err != nil {
		return nil, errors.New(utils.ErrorInvalidBookData)
	}
	b.UUID = uuid

	currentBook, _ := FetchByUUID(uuid, dynaClient)
	if currentBook != nil && len(currentBook.UUID) == 0 {
		return nil, errors.New(utils.ErrorBookDoesNotExists)
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"UUID": &types.AttributeValueMemberS{Value: uuid},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title":  &types.AttributeValueMemberS{Value: b.Title},
			":author": &types.AttributeValueMemberS{Value: b.Author},
			":editor": &types.AttributeValueMemberS{Value: b.Editor},
		},
		TableName:        aws.String(tableName),
		ReturnValues:     types.ReturnValueUpdatedNew,
		UpdateExpression: aws.String("set Title = :title, Author = :author, Editor = :editor"),
	}

	_, err := dynaClient.UpdateItem(context.TODO(), input)
	if err != nil {
		return nil, errors.New(utils.ErrorCouldNotUpdateItem)
	}
	return &b, nil
}

func Delete(uuid string, dynaClient *dynamodb.Client) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"UUID": &types.AttributeValueMemberS{Value: uuid},
		},
	}

	_, err := dynaClient.DeleteItem(context.TODO(), input)
	if err != nil {
		return errors.New(utils.ErrorCouldNotDeleteItem)
	}
	return nil
}
