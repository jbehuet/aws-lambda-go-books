package books

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DeleteBook(req events.APIGatewayProxyRequest, dynaClient *dynamodb.Client, S3Client *s3.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Delete a book
	uuid := req.PathParameters["uuid"]
	if uuid == "" {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(utils.UUIDMissing),
		})
	}

	err := Delete(uuid, dynaClient)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	err = DeleteObject(uuid, S3Client)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return utils.ApiResponse(http.StatusOK, nil)
}
