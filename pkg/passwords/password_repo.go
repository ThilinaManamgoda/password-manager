// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package passwords contains the in memory data manipulations
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

// PasswordEntry struct represents entry in the password DB
type PasswordEntry struct {
	ID       string   `json:"id"`
	Username string    `json:"username"`
	Password string   `json:"password"`
	Labels   []string `json:"labels"`
}

// PasswordDB struct represents password DB
type PasswordDB struct {
	Entries []PasswordEntry `json:"entries"`
}

// PasswordRepository struct handles Password DB
type PasswordRepository struct {
	Encryptor      encrypt.Encryptor
	PasswordFile   utils.PasswordFile
	MasterPassword string
}

func (p *PasswordRepository) loadPasswordDB() ([]byte, error) {
	encryptedData, err := p.PasswordFile.ReadFile()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read password DB file")
	}
	// initialize the Password DB for the first time
	if ! utils.IsValidByteSlice(encryptedData) {
		emptyArray := encryptedData
		return emptyArray, nil
	}
	decryptedData, err := p.Encryptor.Decrypt(encryptedData, p.MasterPassword)
	if err != nil {
		return nil, errors.Wrap(err,"cannot decrypt password DB")
	}
	return decryptedData, nil
}

func (p *PasswordRepository) loadPasswordDBEntries() (*PasswordDB, error) {
	decryptedData, err := p.loadPasswordDB()
	if err != nil {
		return &PasswordDB{}, errors.Wrap(err, "cannot load password DB")
	}

	if ! utils.IsValidByteSlice(decryptedData) {
		return &PasswordDB{}, nil
	}
	var passwordEntries PasswordDB
	if err := json.Unmarshal(decryptedData, &passwordEntries); err != nil {
		return &PasswordDB{}, err
	}
	return &passwordEntries, nil
}

func (p *PasswordRepository) savePasswordDB(passwordDB *PasswordDB, masterPassword string) error {
	passwordDBJSON, err := json.Marshal(passwordDB)
	if err != nil {
		return errors.Wrap(err, "cannot marshal the password DB")
	}
	encryptedData, err := p.Encryptor.Encrypt(passwordDBJSON, masterPassword)
	if err != nil {
		return errors.Wrap(err, "cannot encrypt password DB")
	}
	err = p.PasswordFile.WriteToFile(encryptedData)
	if err != nil {
		return errors.Wrap(err, "cannot write to password DB file")
	}
	return nil
}

// Add method add new password entry to Password DB
func (p *PasswordRepository) Add(id, uN, password string, labels []string) error {
	if id == "" {
		return errors.New("please specify thee ID")
	}
	passwordDB, err := p.loadPasswordDBEntries()
	if err != nil {
		return errors.Wrap(err, "cannot load password DB entries")
	}

	if isIDExists(id, passwordDB.Entries) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}

	passwordDB.Entries = append(passwordDB.Entries, PasswordEntry{
		ID:       id,
		Username: uN,
		Password: password,
		Labels:   labels,
	})
	err = p.savePasswordDB(passwordDB, p.MasterPassword)
	if err != nil {
		return errors.Wrap(err, "cannot save passoword")
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

// GetPassword method retrieve password entry from Password DB
func (p *PasswordRepository) GetPassword(id string, showPassword bool) error {
	passwordDB, err := p.loadPasswordDB()
	if err != nil {
		return err
	}
	if ! utils.IsValidByteSlice(passwordDB) {
		return errors.New("no passwords are available. Add one")
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
