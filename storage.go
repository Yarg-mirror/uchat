package main

import "crypto/rand"

type Storage struct {
	key [32]byte
}

func NewStorage() Storage {
	var key [32]byte
	rand.Read(key[:])

	return Storage{
		key: key,
	}
}
