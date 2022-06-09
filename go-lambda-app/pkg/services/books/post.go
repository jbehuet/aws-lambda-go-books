package books

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateBook(req events.APIGatewayProxyRequest, dynaClient *dynamodb.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Create a book
	result, err := Create(req, dynaClient)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}
	return utils.ApiResponse(http.StatusCreated, result)
}
