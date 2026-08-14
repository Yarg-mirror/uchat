package main

import (
	"errors"
	"os"

	"log"
)

type Config struct {
	identity Identity
	storage  Storage
}

func (c *Config) Load() error {
	os.MkdirAll("config", 0750)

	config, err := os.ReadFile("config/config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Pas de configuration, création.")

			id, err := NewIdentity()
			if err != nil {
				log.Fatalln("Impossible de créer une nouvelle identité")
				return err
			}

			c.identity = id
			c.storage = NewStorage()
			return nil
		}

		// Une autre erreur est survenue
		return err
	}

	Unmarshal(config, nil, nil)
}
