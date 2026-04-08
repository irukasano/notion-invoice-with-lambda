package acceptance_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/irukasano/notion-invoice-with-lambda/internal/s3"
)

func TestIssue4Feature_S3Client(t *testing.T) {
	t.Run("put stores pdf with HLD key format", func(t *testing.T) {
		fake := &fakeObjectAPI{}
		client, err := s3.NewClient("invoice-pdf-stg", fake)
		if err != nil {
			t.Fatalf("NewClient returned error: %v", err)
		}

		key := s3.InvoicePDFKey(s3.InvoicePDFKeyInput{
			BillingDate:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			InvoiceNumber: "S-AB202512-01",
			CustomerName:  "株式会社サンプル / 東京",
			TotalAmount:   120000,
		})

		if err := client.PutObject(
			context.Background(),
			key,
			"application/pdf",
			bytes.NewReader([]byte("%PDF-1.4")),
		); err != nil {
			t.Fatalf("PutObject returned error: %v", err)
		}

		if len(fake.putCalls) != 1 {
			t.Fatalf("put call count mismatch: got %d", len(fake.putCalls))
		}
		got := fake.putCalls[0]
		if got.Bucket != "invoice-pdf-stg" {
			t.Fatalf("bucket mismatch: got %q", got.Bucket)
		}
		wantKey := "2025/12/2025-12-01_S-AB202512-01_株式会社サンプル_東京_イルカシステム_120000.pdf"
		if got.Key != wantKey {
			t.Fatalf("key mismatch: got %q want %q", got.Key, wantKey)
		}
		if got.ContentType != "application/pdf" {
			t.Fatalf("content type mismatch: got %q", got.ContentType)
		}
		if string(got.Body) != "%PDF-1.4" {
			t.Fatalf("body mismatch: got %q", string(got.Body))
		}
	})

	t.Run("get loads object with same bucket and key", func(t *testing.T) {
		fake := &fakeObjectAPI{
			getOutput: s3.GetObjectOutput{
				Body:        io.NopCloser(strings.NewReader("%PDF-1.4")),
				ContentType: "application/pdf",
			},
		}
		client, err := s3.NewClient("invoice-pdf-stg", fake)
		if err != nil {
			t.Fatalf("NewClient returned error: %v", err)
		}

		output, err := client.GetObject(context.Background(), "2025/12/invoice.pdf")
		if err != nil {
			t.Fatalf("GetObject returned error: %v", err)
		}
		defer output.Body.Close()

		body, err := io.ReadAll(output.Body)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if len(fake.getCalls) != 1 {
			t.Fatalf("get call count mismatch: got %d", len(fake.getCalls))
		}
		if fake.getCalls[0].Bucket != "invoice-pdf-stg" {
			t.Fatalf("bucket mismatch: got %q", fake.getCalls[0].Bucket)
		}
		if fake.getCalls[0].Key != "2025/12/invoice.pdf" {
			t.Fatalf("key mismatch: got %q", fake.getCalls[0].Key)
		}
		if string(body) != "%PDF-1.4" {
			t.Fatalf("body mismatch: got %q", string(body))
		}
		if output.ContentType != "application/pdf" {
			t.Fatalf("content type mismatch: got %q", output.ContentType)
		}
	})

	t.Run("put returns backend error", func(t *testing.T) {
		wantErr := errors.New("put failed")
		client, err := s3.NewClient("invoice-pdf-stg", &fakeObjectAPI{putErr: wantErr})
		if err != nil {
			t.Fatalf("NewClient returned error: %v", err)
		}

		err = client.PutObject(
			context.Background(),
			"2025/12/invoice.pdf",
			"application/pdf",
			bytes.NewReader([]byte("%PDF-1.4")),
		)
		if !errors.Is(err, wantErr) {
			t.Fatalf("PutObject error mismatch: got %v want %v", err, wantErr)
		}
	})
}

type fakeObjectAPI struct {
	putCalls  []fakePutCall
	getCalls  []s3.GetObjectInput
	putErr    error
	getErr    error
	getOutput s3.GetObjectOutput
}

type fakePutCall struct {
	Bucket      string
	Key         string
	ContentType string
	Body        []byte
}

func (f *fakeObjectAPI) PutObject(ctx context.Context, input s3.PutObjectInput) error {
	body, err := io.ReadAll(input.Body)
	if err != nil {
		return err
	}

	f.putCalls = append(f.putCalls, fakePutCall{
		Bucket:      input.Bucket,
		Key:         input.Key,
		ContentType: input.ContentType,
		Body:        body,
	})

	return f.putErr
}

func (f *fakeObjectAPI) GetObject(ctx context.Context, input s3.GetObjectInput) (s3.GetObjectOutput, error) {
	f.getCalls = append(f.getCalls, input)
	if f.getErr != nil {
		return s3.GetObjectOutput{}, f.getErr
	}
	return f.getOutput, nil
}
