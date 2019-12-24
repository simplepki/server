package main

import (
	"context"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"strings"
	"errors"
	"net/url"

	"github.com/simplepki/core/keypair"
	"github.com/simplepki/core/types"
	"github.com/simplepki/server/ledger"
	"github.com/simplepki/server/store"
	"github.com/simplepki/server/auth"
)



func HandleRequest(ctx context.Context, event types.CreateCertificateAuthorityEvent) error {
	if event.Token == "" {
		return errors.New("No Auth Token Provided")
	}

	jwtTokenAuth, err := auth.GetJWTAuthorizer("lambda")
	if err != nil {
		return err
	}

	authed, err := jwtTokenAuth.AuthorizeResource(event.Token, "local", event.CAName)
	if err != nil {
		return err
	}

	if !authed {
		return errors.New("Access Denied")
	}

	store := store.AWSSecretsManagerStore{}
	
	uri, err := url.Parse("spiffe://"+event.CAName)
	if err != nil {
		return err
	}

	alreadyExists, err := store.Exists(event.Account, *uri)
	if err != nil {
		return err
	}

	if alreadyExists {
		return errors.New("CA Already Exists")
	}

	kp, err := newCA(ctx, event)
	if err != nil {
		return err
	}

	if len(kp.GetCertificate().URIs) < 1 {
		log.Printf("included uris: %#v\n", kp.GetCertificate().URIs)
		return errors.New("No URI Specified")
	}

	log.Println("building new aurora ledger")
	auroraLedger := ledger.AWSAuroraLedger{}
	log.Println("publishing to aurora ledger")
	err = auroraLedger.Publish(ledger.LedgerRecord{
		Name:        kp.GetCertificate().URIs[0].String(),
		Account:     event.Account,
		Certificate: string(kp.CertificatePEM()),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	store.Put(event.Account, kp)

	return nil
}

func newCA(ctx context.Context, event types.CreateCertificateAuthorityEvent) (keypair.KeyPair, error) {
	var CAName string
	if strings.Contains(event.CAName, "spiffe://") {
		CAName = event.CAName
	} else {
		CAName = fmt.Sprintf("spiffe://%s", event.CAName)
	}

	caName := pkix.Name{
		CommonName: CAName,
	}
	ca := keypair.NewInMemoryKP()

	log.Printf("creating certificate with pkix: %#v\n", caName)
	csr := ca.CreateCSR(caName, []string{})
	caCert := keypair.CsrToCACert(csr)
	ca.Certificate = caCert
	ca.Certificate = ca.IssueCertificate(caCert)
	log.Println("CA generated")

	return ca, nil
}
