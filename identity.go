package main

import (
	"crypto/ed25519"
	"crypto/rand"
)

type Identity struct {
	name    string
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func NewIdentity() (Identity, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return Identity{}, err
	}

	var name [16]byte
	_, err = rand.Read(name[:])
	return Identity{
		name:    string(name[:]),
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}
