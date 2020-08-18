package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Represents uptime monitor result that will be stored in DynamoDB
// Item contains all collected data from single uptime monitor run
type UptimeItem struct {
	RequestID    string `json:"requestId"`    // Uniquely identifies single uptime monitor's run
	UptimeID     string `json:"uptimeId"`
	RunAt        int64  `json:"runAt"`        // Timestamp when the uptime monitor has been invoked
	Host         string `json:"host"`
	StatusCode   int    `json:"statusCode"`
	TTFB         int64  `json:"ttfb"`         // Resulted Time To First Byte in milliseconds
	DNSLookup    int64  `json:"dnslookup"`    // Resulted duration of DNS lookup in milliseconds
	TLSHandshake int64  `json:"tlshandshake"` // Resulted duration of TLS handshake in milliseconds
}

// Store uptime monitor result in DynamoDB table using provide DynamoDB API interface
// Returns error if result cannot be stored in DynamoDB table, otherwise nil
func StoreUptime(uptime UptimeItem, tableName string, db dynamodbiface.DynamoDBAPI) error {
	uptimeItem, err := dynamodbattribute.MarshalMap(uptime)
	if err != nil {
		return err
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      uptimeItem,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	return nil
}
