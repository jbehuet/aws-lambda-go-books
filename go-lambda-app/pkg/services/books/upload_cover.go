package books

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func UploadCover(req events.APIGatewayProxyRequest, s3Client *s3.Client, sqsClient *sqs.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Upload a cover
	err := Upload(req, s3Client)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	// Send message to queue
	messageId, err := SendMessage(req, sqsClient)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}

	return utils.ApiResponse(http.StatusCreated, messageId)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return utils.ApiResponse(http.StatusMethodNotAllowed, utils.ErrorMethodNotAllowed)
}
