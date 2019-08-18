// Copyright © 2019 Thilina Manamgoda
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
	config2 "github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/encrypt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/gojsonq"
)

// PasswordEntry struct represents entry in the password db
type PasswordEntry struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Labels   []string `json:"labels"`
}

// PasswordDB struct represents password db
type PasswordDB struct {
	Entries []PasswordEntry `json:"entries"`
}

// PasswordRepository struct handles Password db
type PasswordRepository struct {
	encryptor     encrypt.Encryptor
	mPassword     string
	rawPasswordDB []byte
	db            *PasswordDB
	file          *fileio.File
}

func loadPasswordDBFile(mPassword string, e encrypt.Encryptor, f *fileio.File) ([]byte, error) {
	encryptedData, err := f.Read()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read password db file")
	}
	// initialize the Password db for the first time
	if !utils.IsValidByteSlice(encryptedData) {
		emptyArray := encryptedData
		return emptyArray, nil
	}
	decryptedData, err := e.Decrypt(encryptedData, mPassword)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decrypt password db")
	}
	return decryptedData, nil
}

func loadPasswordDB(passwordDB []byte) (*PasswordDB, error) {
	var passwordEntries PasswordDB
	if isFirstDBInitialize(passwordDB) {
		return &passwordEntries, nil
	}
	if err := json.Unmarshal(passwordDB, &passwordEntries); err != nil {
		return &PasswordDB{}, errors.Wrapf(err, "cannot unmarshal password db")
	}
	return &passwordEntries, nil
}

func isFirstDBInitialize(db []byte) bool {
	return len(db) == 0
}

func (p *PasswordRepository) savePasswordDB() error {
	passwordDBJSON, err := json.Marshal(p.db)
	if err != nil {
		return errors.Wrap(err, "cannot marshal the password db")
	}
	encryptedData, err := p.encryptor.Encrypt(passwordDBJSON, p.mPassword)
	if err != nil {
		return errors.Wrap(err, "cannot encrypt password db")
	}
	err = p.file.Write(encryptedData)
	if err != nil {
		return errors.Wrap(err, "cannot write to password db file")
	}
	return nil
}

// Add method add new password entry to Password db
func (p *PasswordRepository) Add(id, uN, password string, labels []string) error {
	if id == "" {
		return errors.New("invalid the ID")
	}
	passwordDB := p.db

	if isIDExists(id, passwordDB.Entries) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}

	passwordDB.Entries = append(passwordDB.Entries, PasswordEntry{
		ID:       id,
		Username: uN,
		Password: password,
		Labels:   labels,
	})
	err := p.savePasswordDB()
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

// GetPassword method retrieve password entry from Password db
func (p *PasswordRepository) GetPassword(id string, showPassword bool) error {
	passwordDB := p.rawPasswordDB
	if !utils.IsValidByteSlice(passwordDB) {
		return errors.New("no passwords are available")
	}
	var result []PasswordEntry
	gojsonq.New().JSONString(string(passwordDB)).From("entries").Where("id", "=", id).Out(&result)
	if isResultEmpty(result) {
		return errors.New(fmt.Sprintf("Invalid ID:  %s", id))
	}
	if showPassword {
		fmt.Println(result[0].Password)
	}
	err := clipboard.WriteAll(result[0].Password)
	if err != nil {
		return errors.Wrapf(err, "cannot write to clip board")
	}
	return nil
}

// SearchID will return the password entries if the password ID contains the provide key
func (p *PasswordRepository) SearchID(id string, showPassword bool) ([]PasswordEntry, error) {
	passwordDB := p.rawPasswordDB
	if !utils.IsValidByteSlice(passwordDB) {
		return nil, errors.New("no passwords are available")
	}
	var result []PasswordEntry
	gojsonq.New().JSONString(string(passwordDB)).From("entries").WhereContains("id", id).Out(&result)
	if isResultEmpty(result) {
		return nil, errors.New("cannot find any match")
	}
	return result, nil
}

// SearchLabel will return the password entries if the password labels contains the provide label
func (p *PasswordRepository) SearchLabel(label string, showPassword bool) ([]PasswordEntry, error) {
	if len(p.db.Entries) == 0 {
		return nil, errors.New("no passwords are available")
	}
	var searchResult []PasswordEntry
	for _, e := range p.db.Entries {
		if utils.StringSliceContains(label, e.Labels) {
			searchResult = append(searchResult, e)
		}
	}
	return searchResult, nil
}

// InitPasswordRepo initializes the Password repository
func InitPasswordRepo(mPassword string) (*PasswordRepository, error) {
	config, err := config2.Configuration()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get configuration")
	}
	eFac := &encrypt.Factory{
		ID: config.EncryptorID,
	}
	fSpec := &fileio.File{
		Path: config.PasswordFilePath,
	}
	rawDb, err := loadPasswordDBFile(mPassword, eFac.GetEncryptor(), fSpec)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get Raw PasswordDB")
	}
	db, err := loadPasswordDB(rawDb)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get PasswordDB")
	}
	passwordRepo := &PasswordRepository{
		mPassword:     mPassword,
		encryptor:     eFac.GetEncryptor(),
		rawPasswordDB: rawDb,
		db:            db,
		file:          fSpec,
	}
	return passwordRepo, nil
}
