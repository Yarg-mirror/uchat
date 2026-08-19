package main

import (
	"crypto/ecdh"
)

type Mailbox struct {
	privKey ecdh.PrivateKey
	pubKey  ecdh.PublicKey
}

func NewMailbox() (Mailbox, error) {
	privKey, err := ecdh.X25519().GenerateKey(nil)
	if err != nil {
		return Mailbox{}, err
	}

	return Mailbox{
		privKey: *privKey,
		pubKey:  *privKey.PublicKey(),
	}, nil
}
