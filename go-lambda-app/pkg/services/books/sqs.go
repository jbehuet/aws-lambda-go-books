package books

import (
	"context"
	"errors"
	"jbehuet/aws-lambda-go-books/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const queueName = "book-covers-queue"

func SendMessage(req events.APIGatewayProxyRequest, sqsClient *sqs.Client) (*string, error) {
	// get uuid from query parameters
	uuid := req.PathParameters["uuid"]
	if uuid == "" {
		return nil, errors.New(utils.UUIDMissing)
	}

	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	result, err := sqsClient.GetQueueUrl(context.TODO(), gQInput)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	queueURL := result.QueueUrl

	sMInput := &sqs.SendMessageInput{
		DelaySeconds: 10,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Action": {
				DataType:    aws.String("String"),
				StringValue: aws.String("thumbnail"),
			},
			"UUID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(uuid),
			},
		},
		MessageBody: aws.String("generate thumbnail for " + uuid),
		QueueUrl:    queueURL,
	}

	resp, err := sqsClient.SendMessage(context.TODO(), sMInput)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return resp.MessageId, nil
}
