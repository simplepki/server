package main

import (
	"context"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/simplepki/server/auth"
)

// account, prefix, types string, key []byte, ttlInSeconds int64
type AuthEvent struct {
	Token string `json:"token"`
	TokenType string `json:"token_type"`
	Resource string `json:"resource"`
}

func getJWTKey() ([]byte,error) {
	// Use this code snippet in your app.
	// If you need more information about configurations or implementing the sample code, visit the AWS docs:   
	// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html
	secretName := "jwt"
	//region := "us-west-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return []byte{}, err
	}

	return result.SecretBinary, nil
}

func HandleRequest(ctx context.Context, event AuthEvent) (bool, error) {
	var jwtProvider auth.LocalJWTProvider
	switch event.TokenType {
	case "local":
		jwtProvider = auth.LocalJWTProvider{}
	default:
		//local
		jwtProvider = auth.LocalJWTProvider{}
	}

	key, err := getJWTKey()
	if err != nil {
		return false, err
	}

	validJWT, err := jwtProvider.VerifyJWT(event.Token, key)
	if err != nil {
		return false, err
	}

	return jwtProvider.Authorize(validJWT, event.Resource)
}
