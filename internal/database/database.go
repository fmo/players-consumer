package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Database struct {
	Connection *dynamodb.DynamoDB
}

func NewDbAdapter() *Database {
	return &Database{
		Connection: getConnection(),
	}
}

func getConnection() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}
