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

// Package passwords contains the in memory data manipulations for password repo
package passwords

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/encrypt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

// PasswordEntry struct represents entry in the password db
type PasswordEntry struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrorInvalidID = func(id string) error {
		return errors.New(fmt.Sprintf("Invalid ID:  %s", id))
	}
	ErrorCannotSavePasswordDB = func(err error) error {
		return errors.Wrap(err, "cannot save password")
	}
)

// PasswordDB struct represents password db
type PasswordDB struct {
	Entries map[string]PasswordEntry `json:"entries"`
	Labels  map[string][]string      `json:"labels"`
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
	if _, err := os.Stat(f.Path); err != nil {
		return nil, errors.Wrap(err, "invalid password DB file path")
	}

	encryptedData, err := f.Read()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read password DB file")
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

	if isFirstDBInitialize(passwordDB) {
		return &PasswordDB{
			Entries: map[string]PasswordEntry{},
			Labels:  map[string][]string{},
		}, nil
	}
	var db PasswordDB
	if err := json.Unmarshal(passwordDB, &db); err != nil {
		return &PasswordDB{}, errors.Wrapf(err, "cannot unmarshal password db")
	}
	return &db, nil
}

func isFirstDBInitialize(db []byte) bool {
	return len(db) == 0
}

func (p *PasswordRepository) marshalPasswordDB() ([]byte, error) {
	passwordDBJSON, err := json.Marshal(p.db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal the password db")
	}
	return passwordDBJSON, nil
}

func (p *PasswordRepository) ChangeMasterPassword(newPassword string) error {
	passwordDBJSON, err := p.marshalPasswordDB()
	if err != nil {
		return err
	}
	encryptedData, err := p.encryptor.Encrypt(passwordDBJSON, newPassword)
	if err != nil {
		return errors.Wrap(err, "cannot encrypt password db")
	}
	err = p.file.Write(encryptedData)
	if err != nil {
		return errors.Wrap(err, "cannot write to password db file")
	}
	return nil
}

func (p *PasswordRepository) savePasswordDB() error {
	passwordDBJSON, err := p.marshalPasswordDB()
	if err != nil {
		return err
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

func (p *PasswordRepository) ImportFromCSV(csvFilePath string) error {
	csvfile, err := os.Open(csvFilePath)
	if err != nil {
		return errors.Wrap(err, "Couldn't open the csv file")
	}
	r := csv.NewReader(csvfile)

	first := true
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if first {
			first = false
			continue
		}
		id := record[0]
		uN := record[1]
		password := record[2]
		labels := strings.Split(record[3], ",")

		err = p.addPasswordToRepo(id, uN, password, labels)
		if err != nil {
			return err
		}
	}
	err = p.savePasswordDB()
	if err != nil {
		return ErrorCannotSavePasswordDB(err)
	}
	return nil
}

func (p *PasswordRepository) addPasswordToRepo(id, uN, password string, labels []string) error {
	if p.isIDExists(id) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}
	passwordDBEntries := p.db.Entries
	passwordDBEntries[id] = PasswordEntry{
		ID:       id,
		Username: uN,
		Password: password,
	}
	p.assignLabels(id, labels)
	return nil
}

// Add method add new password entry to Password db
func (p *PasswordRepository) Add(id, uN, password string, labels []string) error {
	if id == "" {
		return errors.New("invalid the ID")
	}
	err := p.addPasswordToRepo(id, uN, password, labels)
	if err != nil {
		return err
	}

	err = p.savePasswordDB()
	if err != nil {
		return errors.Wrap(err, "cannot save passoword")
	}
	return nil
}

func (p *PasswordRepository) isIDExists(id string) bool {
	_, ok := p.db.Entries[id]
	return ok
}

func (p *PasswordRepository) isLabelExists(l string) bool {
	_, ok := p.db.Labels[l]
	return ok
}

// GetPassword method retrieve password entry from Password db
func (p *PasswordRepository) GetPassword(id string, showPassword bool) error {
	passwordDB := p.db.Entries
	if len(passwordDB) == 0 {
		return errors.New("no passwords are available")
	}
	var result PasswordEntry
	result, ok := passwordDB[id]
	if !ok {
		return ErrorInvalidID(id)
	}
	fmt.Println(fmt.Sprintf("Username: %s", result.Username))
	if showPassword {
		fmt.Println(fmt.Sprintf("Password: %s", result.Password))
	} else {
		fmt.Println("Password is copied to the clip board")
		err := clipboard.WriteAll(result.Password)
		if err != nil {
			return errors.Wrapf(err, "cannot write to clip board")
		}
	}
	return nil
}

// SearchID will return the password entries if the password ID contains the provide key
func (p *PasswordRepository) SearchID(id string, showPassword bool) ([]string, error) {
	if p.isDBEmpty() {
		return nil, errors.New("no passwords are available")
	}
	var result []string
	for key := range p.db.Entries {
		if strings.Contains(key, id) {
			result = append(result, key)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("cannot find any match")
	}
	return result, nil
}

func (p *PasswordRepository) isDBEmpty() bool {
	return len(p.db.Entries) == 0
}

// SearchLabel will return the password ids if the password labels contains the provide label
func (p *PasswordRepository) SearchLabel(label string, showPassword bool) ([]string, error) {
	if p.isDBEmpty() {
		return nil, errors.New("no passwords are available")
	}
	var ids []string
	for key, val := range p.db.Labels {
		if strings.Contains(key, label) {
			ids = append(ids, val...)
		}
	}
	uniqueIDs := uniqueStringSlice(ids)
	return uniqueIDs, nil
}

func (p *PasswordRepository) assignLabels(id string, labels []string) {
	for _, val := range labels {
		if p.isLabelExists(val) {
			p.db.Labels[val] = append(p.db.Labels[val], id)
		} else {
			p.db.Labels[val] = []string{id}
		}
	}
}

func (p *PasswordRepository) Remove(id string) error {
	if p.isDBEmpty() {
		return errors.New("no passwords are available")
	}
	if ! p.isIDExists(id) {
		return ErrorInvalidID(id)
	}
	delete(p.db.Entries, id)
	err := p.savePasswordDB()
	if err != nil {
		return ErrorCannotSavePasswordDB(err)
	}
	return nil
}

func uniqueStringSlice(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]struct{})
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}
			u = append(u, val)
		}
	}
	return u
}

// InitPasswordRepo initializes the Password repository
func InitPasswordRepo(mPassword string) (*PasswordRepository, error) {
	conf, err := config.Configuration()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get configuration")
	}
	eFac := &encrypt.Factory{
		ID: conf.EncryptorID,
	}
	fSpec := &fileio.File{
		Path: conf.PasswordDBFilePath,
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
