package main

import (
	"context"
)

type CredentialsEvent struct {
	CAName  string `json:"ca_name"`
	Account string `json:"account"`
}

func HandleRequest(ctx context.Context, event CredentialsEvent) (string, error) {
	
}
