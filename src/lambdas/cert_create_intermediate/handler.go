package main

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"strings"
	"errors"
	"net/url"

	"github.com/simplepki/core/keypair"
	"github.com/simplepki/server/ledger"
	"github.com/simplepki/server/store"
	"github.com/simplepki/server/auth"
)

type CAEvent struct {
	Token string `json:"token"`
	CAName    string `json:"ca_name"`
	InterName string `json:"intermediate_name"`
	Account string `json:"account"`
}

func HandleRequest(ctx context.Context, event CAEvent) error {
	// check not empty
	if event.CAName == "" {
		return errors.New("no ca name specified")
	}
	if event.InterName == "" {
		return errors.New("no intermediate name specified")
	}

	if event.Token == "" {
		return errors.New("No Auth Token Provided")
	}

	jwtTokenAuth, err := auth.GetJWTAuthorizer("lambda")
	if err != nil {
		return err
	}

	authed, err := jwtTokenAuth.AuthorizeResource(event.Token, "local", event.CAName+"/"+event.InterName)
	if err != nil {
		return err
	}

	if !authed {
		return errors.New("Access Denied")
	}

	var CAName, InterName string
	//check for spiffe-ness
	if strings.Contains(event.CAName, "spiffe://") {
		CAName = event.CAName
	} else {
		CAName = fmt.Sprintf("spiffe://%s", event.CAName)
	}

	if strings.Contains(event.InterName, CAName) {
		InterName = event.InterName
	} else {
		InterName = fmt.Sprintf("%s/%s", CAName, event.InterName)
	}


	store := store.AWSSecretsManagerStore{}
	
	caUri, err := url.Parse(CAName)
	if err != nil {
		return err
	}

	caAlreadyExists, err := store.Exists(event.Account, *caUri)
	if err != nil {
		return err
	}

	if !caAlreadyExists {
		return errors.New("CA Doesnt Exist")
	}

	uri, err := url.Parse(InterName)
	if err != nil {
		return err
	}

	alreadyExists, err := store.Exists(event.Account, *uri)
	if err != nil {
		return err
	}

	if alreadyExists {
		return errors.New("Intermediate CA Already Exists")
	}

	// get ca
	log.Println("getting ca with id: ", CAName)
	ca := store.Get(event.Account, *caUri)
	// make new inter
	inter, temp, err := newInter(ctx, InterName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("got intermediate template")
	//sign inter
	inter.Certificate = ca.IssueCertificate(temp)
	log.Println("signed intermediate certificate with CA")
	//publish inter
	auroraLedger := ledger.AWSAuroraLedger{}
	log.Println("publishing to aurora ledger")
	err = auroraLedger.Publish(ledger.LedgerRecord{
		Name:        inter.GetCertificate().URIs[0].String(),
		Account:     event.Account,
		Certificate: string(inter.CertificatePEM()),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	store.Put(event.Account, inter)

	return nil
}

func newInter(ctx context.Context, interName string) (*keypair.InMemoryKP, *x509.Certificate, error) {
	interPkix := pkix.Name{
		CommonName: interName,
	}
	inter := keypair.NewInMemoryKP()

	log.Printf("creating certificate with pkix: %#v\n", interPkix)
	csr := inter.CreateCSR(interPkix, []string{})
	interTemp := keypair.CsrToCert(csr)

	return inter, interTemp, nil
}
