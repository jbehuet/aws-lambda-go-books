package books

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func UpdateBook(req events.APIGatewayProxyRequest, dynaClient *dynamodb.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Update a book
	result, err := Update(req, dynaClient)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}
	return utils.ApiResponse(http.StatusCreated, result)
}
