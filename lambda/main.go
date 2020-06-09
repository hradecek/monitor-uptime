package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodbAPI "github.com/aws/aws-sdk-go/service/dynamodb"
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
	Host       string `json:"host"`
	StatusCode int    `json:"statusCode"` // Resulted status code
	TTFB       int    `json:"ttfb"`       // Measured 'time to first byte'
}

// Get environment variable as string
// If environment variable is not set, then default value is returned instead
func getEnvString(key string, defaultValue string) string {
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
	status, err := uptime.GetUptime(hostUrl, getEnvInt("TIMEOUT", 4))
	if err != nil {
		return nil, err
	}

	return &UptimeMonitorResponse{
		Host:       hostUrl,
		StatusCode: status.StatusCode,
		TTFB:       int(status.TTFB),
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

// Check whether resulted status code matches to requested expectations
func hasExpectedStatusCode(actualStatusCode int, expectedStatusCodes []int) bool {
	for _, expectedStatusCode := range expectedStatusCodes {
		if expectedStatusCode == actualStatusCode {
			return true
		}
	}
	return false
}

// Handle uptime monitor lambda request
// Get uptime response with measured metrics and stored it into DynamoDB
// If resulted status code is not in expected status code provided in request, then send notification to SNS topic
// In case of failure error is returned
func HandleRequest(ctx context.Context, statusReq UptimeMonitorRequest) (UptimeMonitorResponse, error) {
	response, err := response(statusReq.Host)
	if err != nil {
		return UptimeMonitorResponse{}, err
	}

	sessionOptions := session.Options{SharedConfigState: session.SharedConfigEnable}
	db := dynamodbAPI.New(session.Must(session.NewSessionWithOptions(sessionOptions)))
	err = dynamodb.StoreUptime(dynamodb.UptimeItem{
		RequestID:  uuid.New().String(),
		UptimeID:   statusReq.UptimeID,
		RunAt:      time.Now().Unix(),
		Host:       statusReq.Host,
		StatusCode: response.StatusCode,
		TTFB:       response.TTFB,
	}, getEnvString("DYNAMO_TABLE", "uptimes"), db)
	if err != nil {
		return UptimeMonitorResponse{}, err
	}

	if !hasExpectedStatusCode(response.StatusCode, statusReq.StatusCodes) {
		err = sns.PublishUptime(sns.UptimeNotification{
			StatusCode: response.StatusCode,
		}, statusReq.UptimeID, getEnvString("SNS_TOPIC", "uptimes"))
		if err != nil {
			return UptimeMonitorResponse{}, err
		}
	}

	return *response, nil
}

// Main AWS Lambda function
func main() {
	lambda.Start(HandleRequest)
}
