package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
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

func (id *Identity) Sign(message []byte) ([]byte, error) {
	if len(id.privKey) != ed25519.PrivateKeySize {
		err := errors.New("Signing failed")
		return make([]byte, 0), err
	}
	return ed25519.Sign(id.privKey, message), nil
}
