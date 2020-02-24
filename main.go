package main

import (
	"context"
	"os"
	"time"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"github.com/aws/aws-lambda-go/lambda"
)

// Service API
type UptimeRequest struct {
	ClientID string `json:"clientId"`
	UptimeID string `json:"uptimeId"`
	Host string `json:"host"`
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

func Response(host string) (*UptimeResponse, error) {
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

// AWS Lambda API specifcs
func HandleRequest(ctx context.Context, statusReq UptimeRequest) (UptimeResponse, error) {
	response, err := Response(statusReq.Host); if err != nil {
		return UptimeResponse{}, err
	}
	err = storeUptime(UptimeItem{
		RequestID: uuid.New().String(),
		ClientID: statusReq.ClientID,
		UptimeID: statusReq.UptimeID,
		RunAt: time.Now().Unix(),
		Host: statusReq.Host,
		StatusCode: response.StatusCode,
		TTFB: response.TTFB,
	}, getEnvString("DYNAMO_TABLE", "uptimes")); if err != nil {
		return UptimeResponse{}, err
	}
	return *response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
