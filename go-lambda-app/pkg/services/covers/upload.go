package covers

import (
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadCover(req events.APIGatewayProxyRequest, s3Client *s3.Client) (
	*events.APIGatewayProxyResponse,
	error,
) {
	// Upload a cover
	result, err := Upload(req, s3Client)
	if err != nil {
		return utils.ApiResponse(http.StatusBadRequest, utils.ErrorBody{
			ErrorMsg: aws.String(err.Error()),
		})
	}
	return utils.ApiResponse(http.StatusCreated, result)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return utils.ApiResponse(http.StatusMethodNotAllowed, utils.ErrorMethodNotAllowed)
}
