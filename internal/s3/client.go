package s3

import (
	"context"
	"errors"
	"io"
	"strings"
)

var ErrBucketRequired = errors.New("s3 bucket is required")

type ObjectAPI interface {
	PutObject(ctx context.Context, input PutObjectInput) error
	GetObject(ctx context.Context, input GetObjectInput) (GetObjectOutput, error)
}

type PutObjectInput struct {
	Bucket      string
	Key         string
	ContentType string
	Body        io.Reader
}

type GetObjectInput struct {
	Bucket string
	Key    string
}

type GetObjectOutput struct {
	Body        io.ReadCloser
	ContentType string
}

type Client struct {
	bucket string
	api    ObjectAPI
}

func NewClient(bucket string, api ObjectAPI) (*Client, error) {
	trimmedBucket := strings.TrimSpace(bucket)
	if trimmedBucket == "" {
		return nil, ErrBucketRequired
	}
	if api == nil {
		return nil, errors.New("s3 object api is required")
	}

	return &Client{
		bucket: trimmedBucket,
		api:    api,
	}, nil
}

func (c *Client) PutObject(ctx context.Context, key, contentType string, body io.Reader) error {
	return c.api.PutObject(ctx, PutObjectInput{
		Bucket:      c.bucket,
		Key:         key,
		ContentType: contentType,
		Body:        body,
	})
}

func (c *Client) GetObject(ctx context.Context, key string) (GetObjectOutput, error) {
	return c.api.GetObject(ctx, GetObjectInput{
		Bucket: c.bucket,
		Key:    key,
	})
}
