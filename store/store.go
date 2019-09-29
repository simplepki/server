package store

import (
	"github.com/simplepki/core/keypair"
)

type Store interface {
	Exists(account, id string) (bool, error)
	Put(account string, kp keypair.KeyPair) error
	Get(account, id string) keypair.KeyPair
}
