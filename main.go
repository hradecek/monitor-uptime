package main

import (	
	"context"	
	"fmt"
	// "github.com/aws/aws-lambda-go/lambda"
)

// Service API
type PingRequest struct {
	Host string `json:"host"`
}

type PingResponse struct {
	AvgRtt float32 `json:"avgRtt"`
}

// TODO configure via environment variables
func Response(host string) (*PingResponse, error) {
	statistics, err := Ping(host); if err != nil {
		return nil, err
	}

	return &PingResponse{AvgRtt: statistics.AvgRtt}, nil
}

// AWS Lambda API specifcs
func HandleRequest(ctx context.Context, pingReq PingRequest) (PingResponse, error) {
	response, err := Response(pingReq.Host); if err != nil {
		return PingResponse{}, err
	}
	return *response, nil
}

func main() {
	// lambda.Start(HandleRequest)
	fmt.Println(Response("www.google.sk"))
}