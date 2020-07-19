package store

import (
	"github.com/simplepki/core/keypair"
)

type Store interface {
	Init() error
	Exists(KeyPairType keypair.KeyPairType, opts map[string]interface{}) (bool, error)
	Put(kp keypair.KeyPair, opts map[string]interface{}) error
	Get(KeyPairType keypair.KeyPairType, opts map[string]interface{}) keypair.KeyPair
}
