package main

import (
	"context"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(_ context.Context, event events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error){
		log.Printf("Request: %#v\n", event)
		return events.ALBTargetGroupResponse{
			StatusCode: 200,
			StatusDescription: "200 OK",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			IsBase64Encoded: false,
			Body: "Default Page",
			 }, nil
	})
}