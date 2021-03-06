package main

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	//"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/simplepki/core/keypair"
	"github.com/simplepki/core/types"
	"github.com/simplepki/server/store"
	"github.com/simplepki/server/ledger"
	"github.com/simplepki/server/auth"
)

/*
Input Event
{
	"intermediate_name": "test-inter",
	"cert_name":"test-client",
	"account": "test-account",
	"csr":"MIICvDCCAaQCAQAwMjEwMC4GA1UEAxMnc3BpZmZlOi8vdGVzdC1jYS90ZXN0LWludGVyL3Rlc3QtY2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw152HR4jjb3Qz8b6pL0ttMtlEO49JwJ8Ow5Pufn2TfDZw3C0J22NK4T5Y3jUyD9oUBzROeGr+yZUqyQwnKfg2PcCa0KPu3Y2utalVVy2taDq1LD4WR+bSptkeqVpD/suSIQJYoyb0GAHKy1Sr+LuH5Pqr9/QaFm4EoZen/uEPhkxeXpYPc8RkK//bD5XneI2m8qHnFYAzQ+FX/+dhJsfyOZ/qIfCEQFCXFW7w1musCsZcPIPT8ixUzdx8w000Mp3PvWPX9KsEEDm00xANkMpqjy1mNA2VUbpsm78Pn7RTdgU2ep13XFiye3ZknzV6UPocerdMAhNM5G6/TEmNO9UoQIDAQABoEUwQwYJKoZIhvcNAQkOMTYwNDAyBgNVHREEKzAphidzcGlmZmU6Ly90ZXN0LWNhL3Rlc3QtaW50ZXIvdGVzdC1jbGllbnQwDQYJKoZIhvcNAQELBQADggEBAFP2OSOA73fDgNUZDmiRKf1h1mR54FbRfHijd3jqHPEFW4aCBaJTdb+zpplGtO/sd66NY3Pvg29gbIMqWT8gnicv170jZHviSEZmBUF887vc1+H1BG5DLsrLfN1fAV98HvafpdYVsTGf4vR0OyQlRxkpk14/y90KEHPVmIwg8Z4iPgsBj7Ylm8XzMH8lwADRu0LXb84Wo9i8hGw4+M6gM5XPQVw4vaa/b1FWCbvrkNOiSJNGQCG4euav7vcdLEvTrwHVUntL/hOt2ZWjblby/vPP1jvZfXgTyp4CRcNAHabVYlNlfwZc6zr9DGMiYhtIDSsns4qPdwRT5qMHNT6UNWU="
}
*/

func HandleRequest(_ context.Context, event types.SignCertificateEvent) (types.ReturnCertificateEvent, error) {
	if event.InterChain == "" {
		return types.ReturnCertificateEvent{}, errors.New("Missing Intermediate/CA chain")
	}

	if event.Account == "" {
		return types.ReturnCertificateEvent{}, errors.New("Missing Account")
	}

	if len(event.CSR) == 0  {
		return types.ReturnCertificateEvent{}, errors.New("No CSR Provided")
	}

	if event.Token == "" {
		return  types.ReturnCertificateEvent{},errors.New("No Auth Token Provided")
	}

	jwtTokenAuth, err := auth.GetJWTAuthorizer("lambda")
	if err != nil {
		return  types.ReturnCertificateEvent{},err
	}

	authed, err := jwtTokenAuth.AuthorizeResource(event.Token, "local", event.InterChain+"/"+event.CertName)
	if err != nil {
		return  types.ReturnCertificateEvent{},err
	}

	if !authed {
		return  types.ReturnCertificateEvent{},errors.New("Access Denied")
	}

	var InterName string
	if strings.Contains(event.InterChain, "spiffe://") {
		InterName = event.InterChain
	} else {
		InterName = fmt.Sprintf("spiffe://%s", event.InterChain)
	}

	if strings.HasSuffix(InterName, "/") {
		InterName = InterName[0:len(InterName)-1]
	}

	store := store.AWSSecretsManagerStore{}

	interUri, err := url.Parse(InterName)
	if err != nil {
		return types.ReturnCertificateEvent{}, err
	}

	interExists, err := store.Exists(event.Account, *interUri)
	if err != nil {
		return types.ReturnCertificateEvent{}, err
	}

	if !interExists {
		return types.ReturnCertificateEvent{}, errors.New("Intermediate Certificate Doesnt Exist")
	}

	log.Println("getting intermediate with id: ", InterName)
	inter := store.Get(event.Account, *interUri)

	// parse csr
	csrRaw, err := base64.StdEncoding.DecodeString(event.CSR)
	if err != nil {
		log.Println(err.Error())
		return types.ReturnCertificateEvent{}, errors.New("Unable to decode b64 csr")
	}

	// sign certificate
	parsedCSR, err := x509.ParseCertificateRequest(csrRaw)
	if err != nil {
		log.Println(err.Error())
		return types.ReturnCertificateEvent{}, errors.New("Unable to parse csr")
	}

	var CertName string
	if strings.Contains(event.CertName, InterName) {
		CertName = event.CertName
	} else {
		CertName = fmt.Sprintf("%s/%s", InterName, event.CertName)
	}

	pkixName := pkix.Name{
		CommonName: CertName,
	}
	uri, err := url.Parse(pkixName.CommonName)
	if err != nil {
		return types.ReturnCertificateEvent{}, errors.New("Unable to parse uri")
	}

	uriExists := false
	parsedCSR.Subject = pkixName
	if len(parsedCSR.URIs) < 1 {
		parsedCSR.URIs = []*url.URL{uri}
	} else {
		for _, certUri := range parsedCSR.URIs {
			if certUri.String() == uri.String() {
				uriExists = true
				break
			}
		}
	}

	if !uriExists {
		parsedCSR.URIs = []*url.URL{uri}
	}

	certTemp := keypair.CsrToCert(parsedCSR)
	signedCert := inter.IssueCertificate(certTemp)
	log.Printf("client certificate signed: %#v\n", signedCert.Subject.CommonName)
	log.Printf("%#v\n", signedCert)


	auroraLedger := ledger.AWSAuroraLedger{}
	log.Println("publishing to aurora ledger")
	signedPem :=  pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signedCert.Raw,
	})
	err = auroraLedger.Publish(ledger.LedgerRecord{
		Name:        uri.String(),
		Account:     event.Account,
		Certificate: string(signedPem),
	})

	if err != nil {
		return types.ReturnCertificateEvent{}, errors.New("Unable to Publish Certificate")
	}

	chainLedgerRecords, err := auroraLedger.GetChainForRecord(event.Account, *uri)
	if err != nil {
		return types.ReturnCertificateEvent{}, err
	}

	chain := make([]string, len(chainLedgerRecords))
	for idx, record := range chainLedgerRecords {
		chain[idx] = record.Certificate
	}
	
	returnCert := types.ReturnCertificateEvent{
		Cert:  string(signedPem),
		Chain: chain,
	}

	return returnCert, nil
}
