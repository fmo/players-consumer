package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, nil
}

func TestNewS3Service(t *testing.T) {
	service, err := NewS3Service()
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestSave(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test data"))
	}))
	defer ts.Close()

	service := &Service{
		Session: &mockS3Client{},
	}

	err := service.Save("test-key", ts.URL)
	assert.NoError(t, err)
}
