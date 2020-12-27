package sns

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// SNS client mock
type mockSNSClient struct {
	mock.Mock
	snsiface.SNSAPI
}

// SNS erroneous client mock
func (m mockSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

// Given uptime notification,
// When uptime notification is publish into SNS topic for specific uptime monitor,
// Then published MSG contains uptime ID of that uptime monitor as an attribute
func TestPublishUptimeContainsUptimeID(t *testing.T) {
	// Given
	topicARN := "topic-ARN-1"
	uptimeID, _ := uuid.NewRandom()
	uptimeNotification := UptimeNotification{
		Status: STATUS_OK,
	}
	snsClient := mockSNSClient{}
	snsClient.On("Publish", expectedPublishInput(uptimeID, topicARN, uptimeNotification)).Return(&sns.PublishOutput{}, nil)

	// When
	err := PublishUptimeStatus(&uptimeNotification, uptimeID.String(), topicARN, snsClient)

	// Then
	assert.Nil(t, err, "Unexpected error has happened")
	snsClient.AssertExpectations(t)
}

func expectedPublishInput(uptimeID uuid.UUID, expectedTopicARN string, expectedUptimeNotification UptimeNotification) *sns.PublishInput {
	uptimeNotificationJson, _ := json.Marshal(expectedUptimeNotification)

	return &sns.PublishInput{
		TopicArn: aws.String(expectedTopicARN),
		Message:  aws.String(string(uptimeNotificationJson)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"uptimeID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(uptimeID.String()),
			},
		},
	}
}

// Given uptime notification,
// When uptime notification is publish into SNS topic for specific uptime monitor
//      and publish fails,
// Then error is returned.
func TestPublishUptimeSNSFailure(t *testing.T) {
	// Given
	uptimeID, _ := uuid.NewRandom()
	uptimeNotification := UptimeNotification{
		Status: STATUS_OK,
	}
	snsClient := mockSNSClient{}
	snsClient.On("Publish", mock.Anything).Return(&sns.PublishOutput{}, errors.New("cannot publish to SNS topic"))

	// When
	err := PublishUptimeStatus(&uptimeNotification, uptimeID.String(), "topic-ARN-1", snsClient)

	// Then
	assert.NotNil(t, err, "Error was expected to be returned")
	snsClient.AssertExpectations(t)
}
