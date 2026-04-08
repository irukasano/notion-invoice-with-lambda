package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestNewClientRequiresBucket(t *testing.T) {
	_, err := NewClient("", noopObjectAPI{})
	if err == nil {
		t.Fatal("NewClient error = nil")
	}
}

func TestClientPutObjectPassesBucketAndKey(t *testing.T) {
	fake := &captureObjectAPI{}
	client, err := NewClient("invoice-pdf-stg", fake)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	err = client.PutObject(
		context.Background(),
		"2025/12/invoice.pdf",
		"application/pdf",
		bytes.NewReader([]byte("pdf-data")),
	)
	if err != nil {
		t.Fatalf("PutObject returned error: %v", err)
	}

	if fake.lastPut.Bucket != "invoice-pdf-stg" {
		t.Fatalf("bucket mismatch: got %q", fake.lastPut.Bucket)
	}
	if fake.lastPut.Key != "2025/12/invoice.pdf" {
		t.Fatalf("key mismatch: got %q", fake.lastPut.Key)
	}
	if fake.lastPut.ContentType != "application/pdf" {
		t.Fatalf("content type mismatch: got %q", fake.lastPut.ContentType)
	}
	if string(fake.lastPut.Body) != "pdf-data" {
		t.Fatalf("body mismatch: got %q", string(fake.lastPut.Body))
	}
}

func TestClientGetObjectPassesBucketAndReturnsOutput(t *testing.T) {
	fake := &captureObjectAPI{
		getOutput: GetObjectOutput{
			Body:        io.NopCloser(strings.NewReader("pdf-data")),
			ContentType: "application/pdf",
		},
	}
	client, err := NewClient("invoice-pdf-stg", fake)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	output, err := client.GetObject(context.Background(), "2025/12/invoice.pdf")
	if err != nil {
		t.Fatalf("GetObject returned error: %v", err)
	}
	defer output.Body.Close()

	if fake.lastGet.Bucket != "invoice-pdf-stg" {
		t.Fatalf("bucket mismatch: got %q", fake.lastGet.Bucket)
	}
	if fake.lastGet.Key != "2025/12/invoice.pdf" {
		t.Fatalf("key mismatch: got %q", fake.lastGet.Key)
	}

	body, err := io.ReadAll(output.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if string(body) != "pdf-data" {
		t.Fatalf("body mismatch: got %q", string(body))
	}
}

func TestClientPropagatesBackendErrors(t *testing.T) {
	wantErr := errors.New("backend error")
	client, err := NewClient("invoice-pdf-stg", &captureObjectAPI{
		putErr: wantErr,
		getErr: wantErr,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	err = client.PutObject(context.Background(), "invoice.pdf", "application/pdf", bytes.NewReader(nil))
	if !errors.Is(err, wantErr) {
		t.Fatalf("PutObject error mismatch: got %v want %v", err, wantErr)
	}

	_, err = client.GetObject(context.Background(), "invoice.pdf")
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetObject error mismatch: got %v want %v", err, wantErr)
	}
}

type noopObjectAPI struct{}

func (noopObjectAPI) PutObject(ctx context.Context, input PutObjectInput) error {
	return nil
}

func (noopObjectAPI) GetObject(ctx context.Context, input GetObjectInput) (GetObjectOutput, error) {
	return GetObjectOutput{}, nil
}

type captureObjectAPI struct {
	lastPut   capturePutInput
	lastGet   GetObjectInput
	putErr    error
	getErr    error
	getOutput GetObjectOutput
}

type capturePutInput struct {
	Bucket      string
	Key         string
	ContentType string
	Body        []byte
}

func (c *captureObjectAPI) PutObject(ctx context.Context, input PutObjectInput) error {
	body, err := io.ReadAll(input.Body)
	if err != nil {
		return err
	}

	c.lastPut = capturePutInput{
		Bucket:      input.Bucket,
		Key:         input.Key,
		ContentType: input.ContentType,
		Body:        body,
	}

	return c.putErr
}

func (c *captureObjectAPI) GetObject(ctx context.Context, input GetObjectInput) (GetObjectOutput, error) {
	c.lastGet = input
	if c.getErr != nil {
		return GetObjectOutput{}, c.getErr
	}
	return c.getOutput, nil
}
