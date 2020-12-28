package sns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// Represents uptime status. Either OK or FAIL.
type UptimeStatus string

const (
	STATUS_OK   = "OK"
	STATUS_FAIL = "FAIL"
)

// Represents notification sent to SNS topic
type UptimeNotification struct {
	Status UptimeStatus `json:"status"`
}

// Publish uptime notification to SNS topic provided by its ARN
// Published message contains single attribute with uptime ID, which serves for filtering purposes
// Returns error if uptime notification cannot be published to SNS topic, otherwise nil
func PublishUptimeStatus(
	uptimeNotification *UptimeNotification,
	uptimeID string,
	topicARN string,
	snsClient snsiface.SNSAPI) error {
	uptimeNotificationJson, err := json.Marshal(uptimeNotification)
	if err != nil {
		return err
	}

	_, err = snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(string(uptimeNotificationJson)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"uptimeId": {
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
