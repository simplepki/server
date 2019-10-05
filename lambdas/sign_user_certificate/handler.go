package main

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	//"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/simplepki/core/keypair"
	"github.com/simplepki/server/store"
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

type SignEvent struct {
	InterName string `json:"intermediate_name"`
	CertName string `json:"cert_name"`
	Account string `json:"account"`
	CSR string `json:"csr"`
}

type CertEvent struct {
	Cert  string   `json:"cert"`
	Chain []string `json:"chain"`
}

func HandleRequest(_ context.Context, event SignEvent) (CertEvent, error) {
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
	csrRaw, err := base64.StdEncoding.DecodeString(event.CSR)
	if err != nil {
		log.Println(err.Error())
		return CertEvent{}, errors.New("Unable to decode b64 csr")
	}

	// sign certificate
	parsedCSR, err := x509.ParseCertificateRequest(csrRaw)
	if err != nil {
		log.Println(err.Error())
		return CertEvent{}, errors.New("Unable to parse csr")
	}

	var CertName string
	if strings.Contains(event.CertName, InterName) {
		CertName = event.CertName
	} else {
		CertName = InterName + event.CertName
	}

	pkixName := pkix.Name{
		CommonName: CertName,
	}
	uri, err := url.Parse(pkixName.CommonName)
	if err != nil {
		return CertEvent{}, errors.New("Unable to parse uri")
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

	/*
	// get chain
	chain := led.GetChain(csr.Path)

	// return certificate and chain
	returnCert := &CertEvent{
		Cert:  base64.StdEncoding.EncodeToString(signedCert.Raw),
		Chain: chain,
	}

	certBody, err := json.Marshal(returnCert)
	if err != nil {
		log.Fatal(err.Error())
	}
	*/

	return CertEvent{}, nil
}
