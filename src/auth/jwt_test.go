package auth

import (
	"testing"
)

func TestJWT(t *testing.T) {
	prov := LocalJWTProvider{}

	jwt,_ := prov.NewJWT("test", "test/inter", "local", []byte("testkey"),1000)
	t.Log(jwt)

	jwtStruct, err := prov.VerifyJWT(jwt, []byte("testkey"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v\n", jwtStruct)

	notJWT, err := prov.VerifyJWT(jwt, []byte("testkey1"))
	if err == nil {
		t.Fatal("should have errored out")
	}

	t.Log(err, notJWT)
}