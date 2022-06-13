package books

import (
	"bytes"
	"context"
	"errors"
	"io"
	"jbehuet/aws-lambda-go-books/pkg/utils"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucketName = "jbehuet-book-covers"

func StandardHeader(header map[string]string) http.Header {
	h := http.Header{}
	for k, v := range header {
		h.Add(strings.TrimSpace(k), v)
	}
	return h
}

func Upload(req events.APIGatewayProxyRequest, s3Client *s3.Client) error {
	// cf : https://github.com/grokify/go-awslambda/blob/master/multipart.go
	headers := StandardHeader(req.Headers)
	ct := headers.Get("content-type")
	if len(ct) == 0 {
		return errors.New(utils.ContentTypeMissing)
	}

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	if strings.Index(strings.ToLower(strings.TrimSpace(mediatype)), "multipart/") != 0 {
		return errors.New(utils.ContentTypeMissingMultipart)
	}

	paramsInsensitiveKeys := StandardHeader(params)
	boundary := paramsInsensitiveKeys.Get("boundary")
	if len(boundary) == 0 {
		return errors.New(utils.ContentTypeMissingBoundary)
	}

	multipartReader := multipart.NewReader(strings.NewReader(req.Body), boundary)

	part, err := multipartReader.NextPart()
	if err != nil {
		return err
	}
	content, err := io.ReadAll(part)
	if err != nil {
		return err
	}

	// use uuid from query parameters
	uuid := req.PathParameters["uuid"]
	if uuid == "" {
		return errors.New(utils.UUIDMissing)
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uuid),
		Body:   bytes.NewReader(content),
	})

	if err != nil {
		return errors.New(utils.ErrorCouldNotPutObject)
	}

	return nil

}

func GetPresignURL(uuid string, s3Client *s3.Client) (*string, error) {
	s3PresignClient := s3.NewPresignClient(s3Client)

	reponse, err := s3PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uuid),
	})

	if err != nil {
		return nil, errors.New(utils.ErrorCouldNotPresignGetObject)
	}

	return &reponse.URL, err
}

func DeleteObject(uuid string, s3Client *s3.Client) error {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uuid),
	})

	if err != nil {
		return errors.New(utils.ErrorCouldNotDeteleteObject)
	}

	return nil
}
