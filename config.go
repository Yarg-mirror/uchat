package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"

	"log"
)

type Config struct {
	port     uint16
	address  string
	identity Identity
	storage  Storage
}

type configFile struct {
	Config configFileData `json:"config"`
}

type configFileData struct {
	Port     uint16 `json:"port"`
	Address  string `json:"address"`
	Identity struct {
		Name    string `json:"name"`
		PrivKey string `json:"privkey"`
	} `json:"identity"`
	Storage struct {
		Key string `json:"key"`
	} `json:"storage"`
}

func (c *Config) Load() error {
	err := os.Mkdir("config", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Println("Impossible de créer le dossier config.")
		return err
	}

	config, err := os.ReadFile("config/config.json")
	if err != nil && err != os.ErrExist {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Pas de configuration, création.")

			id, err := NewIdentity()
			if err != nil {
				log.Println("Impossible de créer une nouvelle identité.")
				return err
			}

			c.identity = id
			c.storage = NewStorage()
			return nil
		}

		log.Printf("Erreur: %v\n", err)
		return err
	}

	var data configFile
	err = json.Unmarshal(config, &data)
	if err != nil {
		log.Printf("Erreur: %v\n", err)
		return err
	}
	return nil
}

func (c *Config) Save() error {
	var config configFile
	var configData configFileData

	configData.Port = c.port
	configData.Address = c.address
	configData.Identity.Name = c.identity.name
	configData.Identity.PrivKey = base64.StdEncoding.EncodeToString(c.identity.privKey)
	configData.Storage.Key = base64.StdEncoding.EncodeToString(c.storage.key[:])
	config.Config = configData

	output, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("config/config2.json", output, 0640)

	return nil
}

func (c *Config) SetPort(port uint16) {
	c.port = port
}

func (c *Config) SetAddress(address string) {
	c.address = address
}

func (c *Config) SetName(name string) {
	c.identity.name = name
}
