package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"strconv"
)

// Represents uptime monitor result that will be stored in DynamoDB
// Item contains all collected data from single uptime monitor run
type UptimeResultItem struct {
	RequestID    string `json:"requestId"` // Uniquely identifies single uptime monitor's run
	UptimeID     string `json:"uptimeId"`
	RunAt        int64  `json:"runAt"` // Timestamp when the uptime monitor has been invoked
	Host         string `json:"host"`
	StatusCode   int    `json:"statusCode"`
	TTFB         int64  `json:"ttfb"`         // Resulted Time To First Byte in milliseconds
	DNSLookup    int64  `json:"dnslookup"`    // Resulted duration of DNS lookup in milliseconds
	TLSHandshake int64  `json:"tlshandshake"` // Resulted duration of TLS handshake in milliseconds
}

// Store uptime monitor result from single execution in DynamoDB table using provided DynamoDB API interface
// Returns error if result cannot be stored in DynamoDB table, otherwise nil
func StoreUptimeResult(uptime *UptimeResultItem, tableName string, db dynamodbiface.DynamoDBAPI) error {
	if err := putItem(uptime, tableName, db); err != nil {
		return err
	}
	return nil
}

// Put item into Dynamo DB
func putItem(in interface{}, tableName string, db dynamodbiface.DynamoDBAPI) error {
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return err
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	return nil
}

// Store or update uptime's monitor status in DynamoDB table using provided DynamoDB API interface
// For every uptime monitor represented by uptimeID is defined constant threshold and variable failCounter.
// By every call failCounter is incremented. When failCounter cross threshold then true is returned, otherwise false.
// In case of error, non nil error is returned.
func UpdateUptimeStatus(uptimeID string,
	                    threshold string,
	                    tableName string,
	                    db dynamodbiface.DynamoDBAPI) (bool, error) {
	result, err := db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":threshold": {
				N: aws.String(threshold),
			},
			":inc": {
				N: aws.String("1"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"uptimeId": {
				S: aws.String(uptimeID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET threshold=:threshold ADD failCounter :inc"),
	})

	if err != nil {
		return false, err
	}

	return isThresholdCrossed(result), nil
}

// Returns true if failCounter crosses threshold
func isThresholdCrossed(result *dynamodb.UpdateItemOutput) bool {
	failCounter, _ := strconv.Atoi(*result.Attributes["failCounter"].N)
	threshold, _ := strconv.Atoi(*result.Attributes["threshold"].N)
	return failCounter > threshold
}

// If uptime status exists, calling this method clears it (removes from Dynamo DB)
// Returns true if uptime status was cleared
func ClearUptimeStatus(uptimeID string, tableName string, db dynamodbiface.DynamoDBAPI) (bool, error) {
	result, err := db.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"uptimeId": {
				S: aws.String(uptimeID),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return false, err
	}
	return result != nil, nil
}
