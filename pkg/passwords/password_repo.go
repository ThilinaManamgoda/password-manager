/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

// Package password contains the in memory data manipulations
package passwords

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/password-manager/pkg/encrypt"
	"github.com/password-manager/pkg/utils"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/gojsonq"
)

type PasswordEntry struct {
	ID       string   `json:"id"`
	Password string   `json:"password"`
	Labels   []string `json:"labels"`
}

type PasswordDB struct {
	Entries []PasswordEntry `json:"entries"`
}

type PasswordRepository struct {
	Encryptor      encrypt.Encryptor
	PasswordFile   utils.PasswordFile
	MasterPassword string
}

func (p *PasswordRepository) loadPasswordDB() ([]byte, error) {
	encryptedData, err := p.PasswordFile.GetPasswords()
	if err != nil {
		return nil, err
	}
	if ! utils.IsValidByteSlice(encryptedData) {
		emptyArray := encryptedData
		return emptyArray, nil
	}
	decryptedData, err := p.Encryptor.Decrypt(encryptedData, p.MasterPassword)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func (p *PasswordRepository) loadPasswordDBEntries() (*PasswordDB, error) {
	decryptedData, err := p.loadPasswordDB()
	if err != nil {
		return &PasswordDB{}, err
	}

	if ! utils.IsValidByteSlice(decryptedData) {
		return &PasswordDB{}, nil
	} else {
		var passwordEntries PasswordDB
		if err := json.Unmarshal(decryptedData, &passwordEntries); err != nil {
			return &PasswordDB{}, err
		}
		return &passwordEntries, nil
	}
}

func (p *PasswordRepository) savePasswordDB(passwordDB *PasswordDB, masterPassword string) error {
	passwordDBJSON, err := json.Marshal(passwordDB)
	if err != nil {
		return err
	}
	encryptedData, err := p.Encryptor.Encrypt(passwordDBJSON, masterPassword)
	if err != nil {
		return err
	}
	err = p.PasswordFile.StorePasswords(encryptedData)
	if err != nil {
		return err
	}
	return nil
}

func (p *PasswordRepository) Add(id, password string, labels []string) error {
	passwordDB, err := p.loadPasswordDBEntries()
	if err != nil {
		return err
	}

	if isIDExists(id, passwordDB.Entries) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}

	passwordDB.Entries = append(passwordDB.Entries, PasswordEntry{
		ID:       id,
		Password: password,
		Labels:   labels,
	})
	err = p.savePasswordDB(passwordDB, p.MasterPassword)
	if err != nil {
		return err
	}
	return nil
}

func isIDExists(id string, passwordEntries []PasswordEntry) bool {
	for _, val := range passwordEntries {
		if id == val.ID {
			return true
		}
	}
	return false
}

func isResultEmpty(result []PasswordEntry) bool {
	return len(result) == 0
}

func (p *PasswordRepository) GetPassword(id string, showPassword bool) error {
	passwordDB, err := p.loadPasswordDB()
	if err != nil {
		return err
	}
	if ! utils.IsValidByteSlice(passwordDB) {
		return errors.New("No passwords are available. Add one !")
	}
	var result []PasswordEntry
	gojsonq.New().JSONString(string(passwordDB)).From("entries").Where("id", "=", id).Out(&result)
	if isResultEmpty(result) {
		return errors.New(fmt.Sprintf("Invalid ID:  %s", id))
	}
	if showPassword {
		fmt.Println(result[0].Password)
	}
	err = clipboard.WriteAll(result[0].Password)
	if err != nil {
		return err
	}
	return nil
}
