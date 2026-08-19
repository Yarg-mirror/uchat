package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/x509"
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
	mailbox  Mailbox
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
	Mailbox struct {
		PrivKey string `json:"privkey"`
	} `json:"mailbox"`
	Storage struct {
		Key string `json:"key"`
	} `json:"storage"`
}

func (c *Config) Load() error {
	err := os.Mkdir("config", 0750)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Println("Impossible de créer le dossier config.")
		return err
	}

	config, err := os.ReadFile("config/config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Pas de configuration, création.")

			id, err := NewIdentity()
			if err != nil {
				log.Println("Impossible de créer une nouvelle identité.")
				return err
			}
			c.identity = id

			mbox, err := NewMailbox()
			if err != nil {
				log.Println("Impossible de créer une nouvelle boite de réception.")
				return err
			}
			c.mailbox = mbox

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

	c.port = uint16(data.Config.Port)
	c.address = data.Config.Address
	c.identity.name = data.Config.Identity.Name
	c.identity.privKey, _ = base64.StdEncoding.DecodeString(data.Config.Identity.PrivKey)
	c.identity.pubKey = c.identity.privKey.Public().(ed25519.PublicKey)
	decodedPrivKey, err := base64.StdEncoding.DecodeString(data.Config.Mailbox.PrivKey)
	mboxPrivKey, err := x509.ParsePKCS8PrivateKey(decodedPrivKey)
	switch mboxPrivKey.(type) {
	case (*ecdh.PrivateKey):
		c.mailbox.privKey = *mboxPrivKey.(*ecdh.PrivateKey)
	default:
		panic("test")
	}
	c.mailbox.pubKey = *c.mailbox.privKey.PublicKey()
	c.identity.privKey, _ = base64.StdEncoding.DecodeString(data.Config.Identity.PrivKey)
	c.identity.pubKey = c.identity.privKey.Public().(ed25519.PublicKey)
	key, _ := base64.StdEncoding.DecodeString(data.Config.Storage.Key)
	c.storage.key = [32]byte(key)

	return nil
}

func (c *Config) Save() error {
	var config configFile
	var configData configFileData

	configData.Port = c.port
	configData.Address = c.address
	configData.Identity.Name = c.identity.name
	configData.Identity.PrivKey = base64.StdEncoding.EncodeToString(c.identity.privKey)
	mboxPrivKey, _ := x509.MarshalPKCS8PrivateKey(&c.mailbox.privKey)
	configData.Mailbox.PrivKey = base64.StdEncoding.EncodeToString(mboxPrivKey)
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
