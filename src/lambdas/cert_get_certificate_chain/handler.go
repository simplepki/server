package main

import (
	"context"
	"log"
	"fmt"
	"encoding/json"
	"strings"
	"net/url"

	"github.com/simplepki/server/store"
	"github.com/aws/aws-lambda-go/events"
)

type SignEvent struct {
	InterName string `json:"intermediate_name"`
	Account string `json:"account"`
	CSR []byte `json:"csr"`
}

func HandleRequest(_ context.Context, event SignEvent)error {
	if event.InterName == "" {
		log.Fatal("no intermedaite certificatre name")
	}

	if event.Account == "" {
		log.Fatal("no account specified")
	}

	if len(event.CSR) == 0  {
		log.Fatal("zero length csr provided")
	}

	var InterName string
	if strings.Contains(event.InterName, "spiffe://") {
		InterName = event.InterName
	} else {
		InterName = fmt.Sprintf("spiffe://%s", event.InterName)
	}

	store := store.AWSSecretsManagerStore{}

	interUri, err := url.Parse(InterName)
	if err != nil {
		log.Fatal(err.Error())
	}

	interExists, err := store.Exists(event.Account, *interUri)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !interExists {
		log.Fatal("intermediate doesnt exist")
	}

	log.Println("getting intermediate with id: ", InterName)
	inter := store.Get(event.Account, *interUri)

	// parse csr
	

	return nil
}
