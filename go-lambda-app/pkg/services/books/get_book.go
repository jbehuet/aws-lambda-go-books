package books

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetBook(uuid string, dynaClient *dynamodb.Client, s3Client *s3.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Get a book
	result, err := FetchByUUID(uuid, dynaClient)
	if err != nil {
		return utils.ApiResponse(http.StatusNotFound, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	url, err := GetPresignURL(uuid, s3Client)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}
	result.CoverURL = url

	return utils.ApiResponse(http.StatusOK, result)
}
