package dynamodb

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// DynamoDB mock
type mockDynamoDBClient struct {
	failCounter string
	threshold string
	clearedUptimeId string
	dynamodbiface.DynamoDBAPI
}

func (m mockDynamoDBClient) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	if m.clearedUptimeId != "" {
		return &dynamodb.DeleteItemOutput{
			Attributes: map[string]*dynamodb.AttributeValue{
				"uptimeId": {
					S: aws.String(m.clearedUptimeId),
				},
			},
		}, nil
	} else {
		return &dynamodb.DeleteItemOutput{}, nil
	}
}

func (m mockDynamoDBClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, nil
}

func (m mockDynamoDBClient) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{
		Attributes: map[string]*dynamodb.AttributeValue{
			"failCounter": {
				N: aws.String(m.failCounter),
			},
			"threshold": {
				N: aws.String(m.threshold),
			},
		},
	}, nil
}

// DynamoDB erroneous mock
type mockDynamoDBClientBroken struct {
	dynamodbiface.DynamoDBAPI
}

func (m mockDynamoDBClientBroken) DeleteItem(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, errors.New("cannot delete item from dynamodb")
}

func (m mockDynamoDBClientBroken) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, errors.New("cannot put item into dynamodb")
}

func (m mockDynamoDBClientBroken) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return &dynamodb.UpdateItemOutput{}, errors.New("cannot update item in dynamodb")
}

// Given uptime item has been created,
// When uptime item is stored into DynamoDB
//      and DynamoDB PutItem operation fails,
// Then error is returned.
func TestStoreUptimeResultFailure(t *testing.T) {
	// Given
	requestID, _ := uuid.NewRandom()
	uptimeID, _ := uuid.NewRandom()
	uptimeItem := UptimeResultItem{
		RequestID:  requestID.String(),
		UptimeID:   uptimeID.String(),
		RunAt:      time.Now().Unix(),
		Host:       "https://www.hopefuly-not-existing-host.com",
		StatusCode: 200,
		TTFB:       100,
	}

	// When
	err := StoreUptimeResult(&uptimeItem, "SampleTable", mockDynamoDBClientBroken{})

	// Then
	assert.NotNil(t, err, "Error was expected to be returned")
}

// Given uptime item has been create,
// When uptime item is stored into DynamoDB
//      and DynamoDB PutItem operation succeeds,
// Then nil is returned
func TestStoreUptimeResultSuccess(t *testing.T)  {
	// Given
	requestID, _ := uuid.NewRandom()
	uptimeID, _ := uuid.NewRandom()
	uptimeItem := UptimeResultItem{
		RequestID:  requestID.String(),
		UptimeID:   uptimeID.String(),
		RunAt:      time.Now().Unix(),
		Host:       "https://www.hopefuly-not-existing-host.com",
		StatusCode: 200,
		TTFB:       100,
	}

	// When
	err := StoreUptimeResult(&uptimeItem, "SampleTable", mockDynamoDBClient{})

	// Then
	assert.Nil(t, err, "Error was not expected to be returned")
}

// When uptime status is updated
//      and provided threshold is crossed
// Then true is returned
func TestUpdateUptimeStatusSuccessThresholdCrossed(t *testing.T) {
	// When
	res, err := UpdateUptimeStatus("anyUptimeId", "1", "anyTableName", mockDynamoDBClient{
		threshold: "2",
		failCounter: "3",
	})

	// Then
	assert.Nil(t, err, "Error was not expected to be returned")
	assert.True(t, res, "Result was expected to be true")
}

// When uptime status is updated
//      and provided threshold is not crossed
// Then false is returned
func TestUpdateUptimeStatusSuccessThresholdNotCrossed(t *testing.T) {
	// When
	res, err := UpdateUptimeStatus("anyUptimeId", "1", "anyTableName", mockDynamoDBClient{
		threshold: "2",
		failCounter: "2",
	})

	// Then
	assert.Nil(t, err, "Error was not expected to be returned")
	assert.False(t, res, "Result was expected to be false")
}

// When uptime status is updated
//      and error occurs
// Then non-nil error is returned
func TestUpdateUptimeStatusFailure(t *testing.T) {
	// When
	_, err := UpdateUptimeStatus("anyUptimeId", "1", "anyTableName", mockDynamoDBClientBroken{})

	// Then
	assert.NotNil(t, err, "Error was expected to be returned")
}

// Given uptime status exists
// When uptime status is cleared
// Then uptime status is cleared
//      and true is returned
func TestClearUptimeStatusCleared(t *testing.T) {
	// When
	res, err := ClearUptimeStatus("anyUptimeId", "anyTableName", mockDynamoDBClient{
		clearedUptimeId: "anyUptimeId",
	})

	// Then
	assert.Nil(t, err, "Error was not expected to be returned")
	assert.True(t, res, "Result was expected to be true")
}

// Given uptime status does not exist
// When uptime status is cleared
// Then uptime status is cleared
//      and false is returned
func TestClearUptimeStatusNotCleared(t *testing.T) {
	// When
	res, err := ClearUptimeStatus("anyUptimeId", "anyTableName", mockDynamoDBClient{})

	// Then
	assert.Nil(t, err, "Error was not expected to be returned")
	assert.False(t, res, "Result was expected to be false")
}

// When uptime status is cleared
//      and error occurs
// Then non-nil error is returned
func TestClearUptimeStatusFailure(t *testing.T) {
	// When
	_, err := ClearUptimeStatus("anyUptimeId", "anyTableName", mockDynamoDBClientBroken{})

	// Then
	assert.NotNil(t, err, "Error was expected to be returned")
}
