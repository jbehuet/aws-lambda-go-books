package covers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucketName = "jbehuet-book-covers"

type File struct {
	Content       string
	FileName      string
	FileExtension string
}

func StandardHeader(header map[string]string) http.Header {
	h := http.Header{}
	for k, v := range header {
		h.Add(strings.TrimSpace(k), v)
	}
	return h
}

func Upload(req events.APIGatewayProxyRequest, s3Client *s3.Client) (*File, error) {
	// cf : https://github.com/grokify/go-awslambda/blob/master/multipart.go
	headers := StandardHeader(req.Headers)
	ct := headers.Get("content-type")
	if len(ct) == 0 {
		return nil, errors.New("content type missing")
	}

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if strings.Index(strings.ToLower(strings.TrimSpace(mediatype)), "multipart/") != 0 {
		return nil, errors.New("content type missing multipart")
	}

	paramsInsensitiveKeys := StandardHeader(params)
	boundary := paramsInsensitiveKeys.Get("boundary")
	if len(boundary) == 0 {
		return nil, errors.New("content type missing boundary")
	}

	multipartReader := multipart.NewReader(strings.NewReader(req.Body), boundary)

	part, err := multipartReader.NextPart()
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(part)
	if err != nil {
		return nil, err
	}

	// create a unique file name for the file
	// tempFileName := uuid.NewString() + filepath.Ext(part.FileName())

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(part.FileName()),
		Body:   bytes.NewReader(content),
	})

	fmt.Println(err)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	file := File{
		Content:       string(content),
		FileName:      part.FileName(),
		FileExtension: filepath.Ext(part.FileName()),
	}

	return &file, nil

}
