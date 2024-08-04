package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/fmo/players-consumer/config"
	"github.com/fmo/players-consumer/internal/models"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Service struct {
	Session s3iface.S3API
	logger  *log.Logger
}

func NewS3Service(l *log.Logger) (*Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetAwsRegion())},
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		Session: s3.New(sess),
		logger:  l,
	}, nil
}

func (s Service) Save(player models.Player) (imageAlreadyUploaded bool, err error) {
	s3Key := fmt.Sprintf("%s.png", player.Id)

	s3Bucket := config.GetS3Bucket()

	if s.checkImageAlreadyUploaded(s3Bucket, s3Key) == true {
		return true, nil
	}

	// Download the file
	resp, err := http.Get(player.Photo)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Upload the file to S3
	_, err = s.Session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return false, err
	}

	return false, nil
}

func (s Service) checkImageAlreadyUploaded(s3Bucket, s3Key string) bool {
	_, err := s.Session.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
	})

	if err == nil {
		return true
	}

	return false
}
