package auth

import (
	"os"
	"errors"
	"encoding/json"

	"github.com/simplepki/core/types"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type JWTAuthorizer interface {
	AuthorizeResource(jwt string, jwtType string, resource string) (bool, error)
}

type LambdaJWTAuthorizer struct {
	ARN string
}

func GetJWTAuthorizer(authType string) (JWTAuthorizer,error) {
	switch authType {
	case "lambda":
		arn := os.Getenv("JWT_AUTH_ARN")
		if arn == "" {
			return LambdaJWTAuthorizer{}, errors.New("JWT_AUTH_ARN not set")
		}
		return LambdaJWTAuthorizer{ARN: arn}, nil
	default:
		arn := os.Getenv("JWY_AUTH_ARN")
		if arn == "" {
			return LambdaJWTAuthorizer{}, errors.New("JWT_AUTH_ARN not set")
		}
		return LambdaJWTAuthorizer{ARN: arn}, nil
	}
}

func (l LambdaJWTAuthorizer) AuthorizeResource(jwt, jwtType, resource string ) (bool, error) {
	lambdaEvent := types.AuthorizeCredentialsEvent{
		Token: jwt,
		TokenType: jwtType,
		Resource: resource,
	}

	jsonEvent, err := json.Marshal(&lambdaEvent)
	if err != nil {
		return false, err
	}

	lambdaInput := &lambda.InvokeInput{
		FunctionName: aws.String(l.ARN),
		Payload: jsonEvent,
	}

	lambdaSvc := lambda.New(session.New())
	lambdaOutput, err := lambdaSvc.Invoke(lambdaInput)
	if err != nil {
		return false, err
	}

	var result bool
	err = json.Unmarshal(lambdaOutput.Payload, &result)
	if err != nil {
		return false, err
	}
	

	return result, nil
}