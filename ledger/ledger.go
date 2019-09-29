package ledger

type LedgerRecord struct {
	Certificate []byte
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
	GetChainForRecord(string) []LedgerRecord
	GetAllForAccount(string) []LedgerRecordHierarchy
}
