package main

import (	
	"context"	
	"strings"
	// "github.com/aws/aws-lambda-go/lambda"
	"fmt"
)

// Service API
type StatusRequest struct {
	Host string `json:"host"`
}

type StatusResponse struct {
	Host string `json:"host"`
	StatusCode int `json:"statusCode"`
	TTFB int `json:"ttfb"`
}

// TODO configure via environment variables
func Response(host string) (*StatusResponse, error) {
	hostUrl := addProtocol(host)
	status, err := GetStatus(hostUrl); if err != nil {
		return nil, err
	}

	return &StatusResponse{Host: hostUrl,
						   StatusCode: status.StatusCode,
						   TTFB: int(status.TTFB)}, nil
}

func addProtocol(host string) string {
	if !(strings.HasPrefix(host, "http://") && strings.HasPrefix(host, "https://")) {
		return "https://" + host
	}
	return host
}

// AWS Lambda API specifcs
func HandleRequest(ctx context.Context, statusReq StatusRequest) (StatusResponse, error) {
	response, err := Response(statusReq.Host); if err != nil {
		return StatusResponse{}, err
	}
	return *response, nil
}

func main() {
	fmt.Println(HandleRequest(nil, StatusRequest{Host: "www.google.sk"}))
	// lambda.Start(HandleRequest)
}