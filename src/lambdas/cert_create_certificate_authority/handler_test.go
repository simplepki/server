package main

import (
	"context"
	"regexp"
	"testing"
)

func TestCreateCA(t *testing.T) {

	testEvent := CAEvent{
		CAName: "testing",
	}

	kp, err := newCA(context.Background(), testEvent)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("b64: ", kp.Base64Encode())
	matched, _ := regexp.MatchString(`[A-Za-z0-9\+\/]+=?`, kp.Base64Encode())
	if !matched {
		t.Fatal("response not base 64")
	}

	t.Log("size: ", len([]byte(kp.Base64Encode())))
}
