package sns

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/google/uuid"
	"testing"
)

// SNS erroneous mock
type mockSNSClientBroken struct {
	snsiface.SNSAPI
}

func (m mockSNSClientBroken) Publish(*sns.PublishInput) (*sns.PublishOutput, error) {
	return &sns.PublishOutput{}, errors.New("cannot publish into sns topic")
}

// When uptime notification is publish into SNS topic for specific uptime monitor,
// Then published MSG contains uptime ID of that uptime monitor as an attribute
func TestPublishUptimeContainsUptimeID(t *testing.T) {
}

// When uptime notification is publish into SNS topic for specific uptime monitor
//      and publish fails,
// Then error is returned.
func TestPublishUptimeSNSFailure(t *testing.T) {
	uptimeID, _ := uuid.NewRandom()
	uptimeNotification := UptimeNotification{
		StatusCode: 200,
	}

	// When
	err := PublishUptime(uptimeNotification, uptimeID.String(), "topic-ARN-1", mockSNSClientBroken{})

	if err == nil {
		t.Error("Error was expected to be returned")
	}
}
