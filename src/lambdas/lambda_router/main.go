package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lambdasvc "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	lambda.Start(func(_ context.Context, event events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error){
		log.Printf("Request: %#v\n", event)

		var response []byte = []byte{}
		var err error = nil
		switch event.Path {
		case "/create_ca":
			log.Println("executing create_ca lambda")
			response, err = runLambda("CREATE_CA", []byte(event.Body))
		case "/create_intermediate":
			log.Println("executing create_intermediate lambda")
			response, err = runLambda("CREATE_INTERMEDIATE", []byte(event.Body))
		case "/sign_csr":
			log.Println("executing sign_csr lambda")
			response, err = runLambda("SIGN_CSR", []byte(event.Body))
		}

		if err != nil {
			log.Printf("error: %s\n", err.Error())
			return events.ALBTargetGroupResponse{
				StatusCode: 500,
				StatusDescription: "500 Internal Server Error",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				IsBase64Encoded: false,
				Body: `{"error": "contact support"}`,
			}, err
		}

		if len(response) == 0 {
			log.Printf("Empty repsonse\n")
			return events.ALBTargetGroupResponse{
				StatusCode: 500,
				StatusDescription: "500 Internal Server Error",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				IsBase64Encoded: false,
				Body: `{"error": "no response"}`,
			}, err
		}

		// default return
		return events.ALBTargetGroupResponse{
			StatusCode: 200,
			StatusDescription: "200 OK",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			IsBase64Encoded: false,
			Body: string(response),
			 }, nil
	})
}

func runLambda(envVar string, jsonEvent []byte) ([]byte, error) {
	lambdaArn := os.Getenv(envVar)

	lambdaInput := &lambdasvc.InvokeInput{
		FunctionName: aws.String(lambdaArn),
		Payload: jsonEvent,
	}

	svc := lambdasvc.New(session.New())
	lambdaOutput, err := svc.Invoke(lambdaInput)
	if err != nil {
		return []byte{}, err
	}

	return lambdaOutput.Payload, nil
}