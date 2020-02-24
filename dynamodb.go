package main

import (
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UptimeItem struct {
	RequestID string `json:"requestId"`
	ClientID string `json:"clientId"`
	UptimeID string `json:"uptimeId"`
	RunAt int64 `json:"runAt"`
	Host string `json:"host"`
	StatusCode int `json:"statusCode"`
	TTFB int `json:"ttfb"`
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
