package ledger

import (
	"net/url"
)

type LedgerRecord struct {
	Certificate string
	Name        string
	Account     string
}

type LedgerRecordHierarchy struct {
	Root     LedgerRecord
	Children []LedgerRecordHierarchy
}

type Ledger interface {
	Publish(LedgerRecord)
	GetRecordByName(string) LedgerRecord
	GetRecordByAccount(string) []LedgerRecord
	GetChainForRecord(url.URL) ([]LedgerRecord, error)
	GetAllForAccount(string) []LedgerRecordHierarchy
}
