package sns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// Represents notification sent to SNS topic
type UptimeNotification struct {
	StatusCode int `json:"status"`
}

// Publish uptime notification to SNS topic provided by its ARN
// Published message contains single attribute with uptime ID, which serves for filtering purposes
// Returns error if uptime notification cannot be published to SNS topic, otherwise nil
func PublishUptime(uptimeNotification UptimeNotification, uptimeID string, topicARN string, snsClient snsiface.SNSAPI) error {
	uptimeNotificationJson, err := json.Marshal(uptimeNotification)
	if err != nil {
		return err
	}

	_, err = snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(string(uptimeNotificationJson)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"uptimeID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(uptimeID),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
