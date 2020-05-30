package main

import (
	"context"
	"os"
	"time"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

// Service API
type UptimeRequest struct {
	UptimeID string `json:"uptimeId"`
	Host string `json:"host"`
	StatusCodes []int `json:"statusCodes"`
}

type UptimeResponse struct {
	Host string `json:"host"`
	StatusCode int `json:"statusCode"`
	TTFB int `json:"ttfb"`
}

func getEnvString(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, _ := strconv.Atoi(value)
		return intValue
	}
	return defaultValue
}

func response(host string) (*UptimeResponse, error) {
	hostUrl := addProtocol(host)
	status, err := GetUptime(hostUrl, getEnvInt("TIMEOUT", 4)); if err != nil {
		return nil, err
	}

	return &UptimeResponse{Host: hostUrl,
						   StatusCode: status.StatusCode,
						   TTFB: int(status.TTFB)}, nil
}

func addProtocol(host string) string {
	if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
		return "https://" + host
	}
	return host
}

func hasExpectedStatusCode(actualStatusCode int, expectedStatusCodes []int) bool {
	for _, expectedStatusCode := range expectedStatusCodes {
		if expectedStatusCode == actualStatusCode {
			return true
		}
	}
	return false
}

// AWS Lambda API specifcs
func HandleRequest(ctx context.Context, statusReq UptimeRequest) (UptimeResponse, error) {
	response, err := response(statusReq.Host); if err != nil {
		return UptimeResponse{}, err
	}

	db := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable})))
	err = StoreUptime(UptimeItem{
		RequestID: uuid.New().String(),
		UptimeID: statusReq.UptimeID,
		RunAt: time.Now().Unix(),
		Host: statusReq.Host,
		StatusCode: response.StatusCode,
		TTFB: response.TTFB,
	}, getEnvString("DYNAMO_TABLE", "uptimes"), db); if err != nil {
		return UptimeResponse{}, err
	}

	if !hasExpectedStatusCode(response.StatusCode, statusReq.StatusCodes) {
		err = PublishUptime(UptimeNotification{
			StatusCode: response.StatusCode,
		}, statusReq.UptimeID, getEnvString("SNS_TOPIC", "sns_topic")); if err != nil {
			return UptimeResponse{}, err
		}
	}

	return *response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
