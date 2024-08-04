package config

import (
	"fmt"
	"os"
)

func GetS3Bucket() string {
	return os.Getenv("S3_BUCKET")
}

func GetAwsRegion() string {
	return os.Getenv("AWS_REGION")
}

func GetApiPort() string {
	return fmt.Sprintf(":%s", os.Getenv("API_PORT"))
}
