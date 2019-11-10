package store

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/simplepki/core/keypair"
)

type AWSSecretsManagerStore struct{}

func (sms AWSSecretsManagerStore) Exists(account string, id url.URL) (bool, error) {
	nameString := fmt.Sprintf("%s/%s%s", account, id.Host, id.Path)

	config := &aws.Config{}
	session := session.Must(session.NewSession(config))

	sm := secretsmanager.New(session)

	describeSecretInput := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(nameString),
	}

	describeResponse, err := sm.DescribeSecret(describeSecretInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}

	if describeResponse == nil {
		return false, nil
	}

	return true, nil
}

func (sms AWSSecretsManagerStore) Put(account string, kp keypair.KeyPair) error {

	if len(kp.GetCertificate().URIs) < 1 {
		return errors.New("No URI supplied")
	}

	url := kp.GetCertificate().URIs[0]
	nameString := fmt.Sprintf("%s/%s%s", account, url.Host, url.Path)

	config := &aws.Config{}
	session := session.Must(session.NewSession(config))

	sm := secretsmanager.New(session)

	secretsInput := &secretsmanager.CreateSecretInput{
		Description:  aws.String("CA KeyPair"),
		Name:         aws.String("cert@"+nameString),
		SecretString: aws.String(kp.Base64Encode()),
	}

	resp, err := sm.CreateSecret(secretsInput)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("create secret response: %#v\n", resp)
	return nil
}

func (sms AWSSecretsManagerStore) Get(account string, id url.URL) keypair.KeyPair {
	nameString := fmt.Sprintf("%s/%s%s", account, id.Host, id.Path)

	config := &aws.Config{}
	session := session.Must(session.NewSession(config))

	sm := secretsmanager.New(session)

	query := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("cert@"+nameString),
	}

	output, err := sm.GetSecretValue(query)
	if err != nil {
		log.Fatal(err)
	}

	kp := &keypair.InMemoryKP{}
	kp.Base64Decode(*output.SecretString)

	return kp
}
