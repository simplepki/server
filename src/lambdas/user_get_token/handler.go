package main

import (
	"context"
	"crypto/rand"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	//"encoding/base64"
	//"fmt"

	"github.com/simplepki/server/auth"
	"github.com/simplepki/core/types"
)

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
		if awsError, ok := err.(awserr.Error); ok {
			switch awsError.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				// create a random key
				key := make([]byte, 256)
				_, err := rand.Read(key)
				if err != nil {
					return []byte{}, err
				}
				
				createInput := &secretsmanager.CreateSecretInput{
					Name: aws.String(secretName),
					SecretBinary: key,
				}

				_, err = svc.CreateSecret(createInput)
				if err != nil {
					return []byte{}, err
				}

				return key, nil
			default:
				return []byte{}, err
			}
		}
	}

	return result.SecretBinary, nil
}

func HandleRequest(ctx context.Context, event types.CreateCredentialsEvent) (string, error) {
	var jwtProvider auth.LocalJWTProvider
	switch event.Type {
	case "local":
		jwtProvider = auth.LocalJWTProvider{}

	default:
		//local
		jwtProvider = auth.LocalJWTProvider{}
	}

	key, err := getJWTKey()
	if err != nil {
		return "", err
	}

	return jwtProvider.NewJWT(event.Account, event.Prefix, event.Type, key, event.TTL)
}
