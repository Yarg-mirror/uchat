package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"os"

	"log"
)

type ConfigFile struct {
	Config ConfigFileData `json:"config"`
}

type ConfigFileData struct {
	Port     uint16 `json:"port"`
	Address  string `json:"address"`
	Identity struct {
		Name    string `json:"name"`
		PrivKey string `json:"privkey"`
	} `json:"identity"`
	Mailbox struct {
		PrivKey string `json:"privkey"`
	} `json:"mailbox"`
	Storage struct {
		Key string `json:"key"`
	} `json:"storage"`
}

func (c *ConfigFile) Load(path string) (*ConfigFileData, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Erreur: %v\n", err)
		return &ConfigFileData{}, err
	}

	var data ConfigFile
	err = json.Unmarshal(configFile, &data)
	if err != nil {
		log.Printf("Erreur: %v\n", err)
		return &ConfigFileData{}, err
	}

	return &data.Config, nil
}

func (c *ConfigFile) Save(path string, config Config) error {
	var configData ConfigFileData

	mboxPrivKey, err := x509.MarshalPKCS8PrivateKey(&config.mailbox.privKey)
	if err != nil {
		log.Printf("Erreur: %v\n", err)
		return err
	}
	configData.Mailbox.PrivKey = base64.StdEncoding.EncodeToString(mboxPrivKey)

	configData.Port = config.port
	configData.Address = config.address
	configData.Identity.Name = config.identity.name
	configData.Identity.PrivKey = base64.StdEncoding.EncodeToString(config.identity.privKey)
	configData.Storage.Key = base64.StdEncoding.EncodeToString(config.storage.key[:])
	c.Config = configData

	output, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(path, output, 0600)

	return nil
}
