package main

import "crypto/ecdh"

type Message struct {
	from      ecdh.PrivateKey
	to        ecdh.PublicKey
	message   []byte
	signature []byte
}
