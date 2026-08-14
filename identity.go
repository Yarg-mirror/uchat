package main

import (
	"crypto/ed25519"
)

type Identity struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func NewIdentity() (Identity, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return Identity{}, err
	}

	return Identity{
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}
