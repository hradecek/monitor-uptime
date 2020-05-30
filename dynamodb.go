package main

import (
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UptimeItem struct {
	RequestID string `json:"requestId"`
	UptimeID string `json:"uptimeId"`
	RunAt int64 `json:"runAt"`
	Host string `json:"host"`
	StatusCode int `json:"statusCode"`
	TTFB int `json:"ttfb"`
}

func StoreUptime(uptime UptimeItem, tableName string, db dynamodbiface.DynamoDBAPI) error {
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
