package ledger

import (
	"testing"
)

func TestSpiffeRegex(t *testing.T) {
	path := "spiffe://examle.com/web/frontend/web1"
	dynamo := AWSDynamoLedger{}
	dynamo.GetChain(path)
}
