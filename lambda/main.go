package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodbAPI "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	snsAPI "github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
	"monitor-uptime/internal/dynamodb"
	"monitor-uptime/internal/sns"
	"monitor-uptime/internal/uptime"
	"os"
	"strconv"
	"strings"
	"time"
)

// Represents uptime monitor service request
type UptimeMonitorRequest struct {
	UptimeID    string `json:"uptimeId"`    // Uptime ID that invoked service
	Host        string `json:"host"`        // Host for which uptime will be invoked
	StatusCodes []int  `json:"statusCodes"` // Expected status code
}

// Represents uptime monitor service response
type UptimeMonitorResponse struct {
	Host         string `json:"host"`
	StatusCode   int    `json:"statusCode"`   // Resulted status code
	TTFB         int64  `json:"ttfb"`         // Measured Time To First Byte in milliseconds
	DNSLookup    int64  `json:"dnslookup"`    // Measured duration of DNS lookup in milliseconds
	TLSHandshake int64  `json:"tlshandshake"` // Measured duration of TLS handshake in milliseconds
}

// Get environment variable as string
// If environment variable is not set, then nil is returned instead
func getEnvString(key string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}
	return nil
}

// Get environment variable as string
// If environment variable is not set, then default value is returned instead
func getEnvStringWithDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// Get environment variable as string
// If environment variable is not set, then default value is returned instead
func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, _ := strconv.Atoi(value)
		return intValue
	}
	return defaultValue
}

// Get uptime monitor response for provided host
// In case of failure error is returned
func response(host string) (*UptimeMonitorResponse, error) {
	hostUrl := sanityHTTPProtocol(host)
	response, err := uptime.GetUptime(hostUrl, getEnvInt("TIMEOUT", 4))
	if err != nil {
		return nil, err
	}

	return &UptimeMonitorResponse{
		Host:         hostUrl,
		StatusCode:   response.StatusCode,
		TTFB:         response.TTFB.Milliseconds(),
		DNSLookup:    response.DNSLookup.Milliseconds(),
		TLSHandshake: response.TLSHandshake.Milliseconds(),
	}, nil
}

// Makes sure that host always contains protocol part
// If not provided explicitly HTTPS is added by default
func sanityHTTPProtocol(host string) string {
	if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
		return "https://" + host
	}
	return host
}

// Stores uptime into Dynamo DB
func storeUptime(statusReq *UptimeMonitorRequest, response *UptimeMonitorResponse, db dynamodbiface.DynamoDBAPI) error {
	tableName := getEnvString("DYNAMO_TABLE_EXECUTIONS")
	if tableName == nil {
		return nil
	}
	return dynamodb.StoreUptimeResult(&dynamodb.UptimeResultItem{
		RequestID:    uuid.New().String(),
		UptimeID:     statusReq.UptimeID,
		RunAt:        time.Now().Unix(),
		Host:         statusReq.Host,
		StatusCode:   response.StatusCode,
		TTFB:         response.TTFB,
		DNSLookup:    response.DNSLookup,
		TLSHandshake: response.TLSHandshake,
	}, *tableName, db)
}

// Updates uptime status in Dynamo DB
// If status has been changed (e.g. cross threshold or uptime went from Fail to OK), then new status is returned
func updateUptimeStatus(statusReq *UptimeMonitorRequest,
	response *UptimeMonitorResponse,
	db dynamodbiface.DynamoDBAPI) (*sns.UptimeStatus, error) {
	var err error
	var notify bool
	var status sns.UptimeStatus

	dbStatus := getEnvStringWithDefault("DYNAMO_TABLE_STATUS", "uptimeStatus")

	if hasExpectedStatusCode(response.StatusCode, statusReq.StatusCodes) {
		status = sns.STATUS_OK
		notify, err = dynamodb.ClearUptimeStatus(statusReq.UptimeID, dbStatus, db)
	} else {
		status = sns.STATUS_FAIL
		notify, err = dynamodb.UpdateUptimeStatus(statusReq.UptimeID, "3", dbStatus, db)
	}

	if err == nil {
		if notify {
			return &status, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

// Notify uptime monitor status via SNS
func notifyUptimeStatus(uptimeId string, status sns.UptimeStatus, sessionOptions *session.Options) error {
	snsTopicName := getEnvString("SNS_TOPIC")

	if snsTopicName != nil {
		snsClient := snsAPI.New(session.Must(session.NewSessionWithOptions(*sessionOptions)))
		return sns.PublishUptimeStatus(&sns.UptimeNotification{Status: status}, uptimeId, *snsTopicName, snsClient)
	}
	return nil
}

// Checks whether resulted status code matches to requested expectations
func hasExpectedStatusCode(actualStatusCode int, expectedStatusCodes []int) bool {
	for _, expectedStatusCode := range expectedStatusCodes {
		if expectedStatusCode == actualStatusCode {
			return true
		}
	}
	return false
}

// Handles uptime monitor lambda request
// Get uptime response with measured metrics and stored it into DynamoDB
// If resulted status code is not in expected status code provided in request, then send notification to SNS topic
// In case of failure error is returned
func HandleRequest(ctx context.Context, req UptimeMonitorRequest) (UptimeMonitorResponse, error) {
	res, err := response(req.Host)
	if err != nil {
		return UptimeMonitorResponse{}, err
	}

	sessionOptions := session.Options{SharedConfigState: session.SharedConfigEnable}
	db := dynamodbAPI.New(session.Must(session.NewSessionWithOptions(sessionOptions)))
	if err = storeUptime(&req, res, db); err != nil {
		return UptimeMonitorResponse{}, err
	}

	var status *sns.UptimeStatus
	status, err = updateUptimeStatus(&req, res, db)
	if err != nil {
		return UptimeMonitorResponse{}, err
	}
	if status != nil {
		if err = notifyUptimeStatus(req.UptimeID, *status, &sessionOptions); err != nil {
			return UptimeMonitorResponse{}, err
		}
	}

	return *res, nil
}

// Main AWS Lambda function
func main() {
	lambda.Start(HandleRequest)
}
