package sns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Represents notification sent to SNS topic
type UptimeNotification struct {
	StatusCode int `json:"status"`
}

// Publish uptime notification to SNS topic provided by its ARN
// Published message contains single attribute with uptime ID, which serves for filtering purposes
// Returns error if uptime notification cannot be published to SNS topic, otherwise nil
func PublishUptime(uptimeNotification UptimeNotification, uptimeID string, topicARN string) error {
	sessionOptions := session.Options{SharedConfigState: session.SharedConfigEnable}
	client := sns.New(session.Must(session.NewSessionWithOptions(sessionOptions)))

	uptimeNotificationJson, err := json.Marshal(uptimeNotification)
	if err != nil {
		return err
	}

	_, err = client.Publish(&sns.PublishInput{
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
