package main


import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type UptimeNotification struct {
	StatusCode int `json:"status"`
}

func PublishUptime(uptime UptimeNotification, uptimeId string, topicArn string) error {
	client := sns.New(session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})))
	uptimeNotification, err := json.Marshal(uptime)
	if err != nil {
		return err
	}

	_, err = client.Publish(&sns.PublishInput{
        TopicArn: aws.String(topicArn),
		Message: aws.String(string(uptimeNotification)),
        MessageAttributes:  map[string]*sns.MessageAttributeValue{
			"uptimeId": {
				DataType: aws.String("String"),
				StringValue: aws.String(uptimeId),
			},
		},
	})

    if err != nil {
		return err
	}

	return nil
}
