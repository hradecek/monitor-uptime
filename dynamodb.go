package main

import (
	"time"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UptimeItem struct {
	RequestID string
	ClientID string
	Timestamp time.Time
	Host string
	StatusCode int
	TTFB int 
}

func storeUptime(uptime UptimeItem, tableName string) error {
	db := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})))
	uptimeItem, err := dynamodbattribute.MarshalMap(uptime)
	if err != nil {
		return err
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item: uptimeItem,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return err
	}

	return nil
}
