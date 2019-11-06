package auth

import (
	"testing"
)

func TestJWT(t *testing.T) {
	prov := LocalJWTProvider{}

	jwt, _ := prov.NewJWT("test", "test/inter", "local", []byte("testkey"), 1000)
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

func TestJWTAuthorization(t *testing.T) {
	prov := LocalJWTProvider{}
	jwt := JWT{
		Prefix: "*",
	}

	resources := []string{"", "test-ca", "test-ca/test-inter"}
	for _, r := range resources {
		result, err := prov.Authorize(jwt, r)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s should match against %s\n", jwt.Prefix, r)
		if result == false {
			t.Fatalf("%s does not match against %s\n", jwt.Prefix, r)
		}
	}

	jwt = JWT{
		Prefix: "test-ca/*",
	}
	resources = []string{"test-ca", "test-ca/inter-1", "test-ca2/inter-1"}
	resourceResult := []bool{false, true, false}
	for idx, r := range resources {
		result, err := prov.Authorize(jwt, r)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("pattern %#v should evaluate %s as %#v\n", jwt.Prefix, r, resourceResult[idx])
		if result != resourceResult[idx] {
			t.Fatalf("pattern %#v should evaluate %s as %#v\n", jwt.Prefix, r, resourceResult[idx])
		}
	}

}
