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
	"encoding/json"
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/encrypt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// Entry struct represents entry in the password db
type Entry struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrInvalidID = func(id string) error {
		return errors.New(fmt.Sprintf("Invalid ID:  %s", id))
	}
	ErrCannotSavePasswordDB = func(err error) error {
		return errors.Wrap(err, "cannot save password")
	}
	ErrNoPasswords          = errors.New("no passwords are available")
	ErrCannotFindMatchForID = func(id string) error {
		return errors.New(fmt.Sprintf("cannot find any match for id %s", id))
	}
)

// DB struct represents password db
type DB struct {
	Entries map[string]Entry    `json:"entries"`
	Labels  map[string][]string `json:"labels"`
}

// Repository struct handles Password db
type Repository struct {
	encryptor encrypt.Encryptor
	mPassword string
	db        *DB
	file      *fileio.File
}

func isPasswordRepoAlreadyInitialized(repoData []byte) bool {
	return utils.IsValidByteSlice(repoData)
}

func loadDBFile(mPassword string, e encrypt.Encryptor, f *fileio.File) ([]byte, error) {
	if exists, err := utils.IsFileExists(f.Path); err != nil {
		return nil, errors.Wrap(err, "cannot load the password DB file")
	} else {
		if !exists {
			return nil, errors.Wrap(err, "cannot find the password DB file")
		}
	}
	encryptedData, err := f.Read()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read password DB file")
	}

	if !isPasswordRepoAlreadyInitialized(encryptedData) {
		return nil, errors.New("password repository is not initialized")
	}

	decryptedData, err := e.Decrypt(encryptedData, mPassword)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decrypt password db")
	}
	return decryptedData, nil
}

func loadDB(passwordDB []byte) (*DB, error) {
	var db DB
	if err := json.Unmarshal(passwordDB, &db); err != nil {
		return &DB{}, errors.Wrapf(err, "cannot unmarshal password db")
	}
	return &db, nil
}

func (p *Repository) marshalDB() ([]byte, error) {
	passwordDBJSON, err := utils.MarshalData(p.db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal the password db")
	}
	return passwordDBJSON, nil
}

func (p *Repository) ChangeMasterPassword(newPassword string) error {
	passwordDBJSON, err := p.marshalDB()
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

func (p *Repository) saveDB() error {
	passwordDBJSON, err := p.marshalDB()
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

func (p *Repository) addPasswordEntryToRepo(id, uN, password string, labels []string) error {
	if p.isIDExists(id) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}
	entries := p.db.Entries
	entries[id] = Entry{
		ID:       id,
		Username: uN,
		Password: password,
	}
	p.assignLabels(id, labels)
	return nil
}

// Add method add new password entry to Password db
func (p *Repository) Add(id, uN, password string, labels []string) error {
	if id == "" {
		return errors.New("invalid the ID")
	}
	err := p.addPasswordEntryToRepo(id, uN, password, labels)
	if err != nil {
		return err
	}

	err = p.saveDB()
	if err != nil {
		return ErrCannotSavePasswordDB(err)
	}
	return nil
}

func (p *Repository) isIDExists(id string) bool {
	_, ok := p.db.Entries[id]
	return ok
}

func (p *Repository) isLabelExists(l string) bool {
	_, ok := p.db.Labels[l]
	return ok
}

// GetPassword method retrieve password entry from Password db
func (p *Repository) GetPassword(id string, showPassword bool) error {
	passwordEntry, err := p.GetPasswordEntry(id)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Username: %s", passwordEntry.Username))
	if showPassword {
		fmt.Println(fmt.Sprintf("Password: %s", passwordEntry.Password))
	} else {
		fmt.Println("Password is copied to the clip board")
		err := clipboard.WriteAll(passwordEntry.Password)
		if err != nil {
			return errors.Wrapf(err, "cannot write to clip board")
		}
	}
	return nil
}

func (p *Repository) GetPasswordEntry(id string) (Entry, error) {
	passwordDB := p.db.Entries
	if len(passwordDB) == 0 {
		return Entry{}, ErrNoPasswords
	}
	var result Entry
	result, ok := passwordDB[id]
	if !ok {
		return Entry{}, ErrInvalidID(id)
	}
	return result, nil
}

func (p *Repository) ChangePasswordEntry(id string, entry Entry) error {
	passwordDB := p.db.Entries
	if len(passwordDB) == 0 {
		return ErrNoPasswords
	}
	passwordDB[id] = entry
	err := p.saveDB()
	if err != nil {
		return ErrCannotSavePasswordDB(err)
	}
	return nil
}

// SearchID will return the password entries if the password ID contains the provide key
func (p *Repository) SearchID(id string, showPassword bool) ([]string, error) {
	if p.isDBEmpty() {
		return nil, ErrNoPasswords
	}
	var result []string
	for key := range p.db.Entries {
		if strings.Contains(key, id) {
			result = append(result, key)
		}
	}
	if len(result) == 0 {
		return nil, ErrCannotFindMatchForID(id)
	}
	return result, nil
}

func (p *Repository) isDBEmpty() bool {
	return len(p.db.Entries) == 0
}

// SearchLabel will return the password ids if the password labels contains the provide label
func (p *Repository) SearchLabel(label string, showPassword bool) ([]string, error) {
	if p.isDBEmpty() {
		return nil, ErrNoPasswords
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

// searchLabelsForID will return the list of labels for given ID
func (p *Repository) searchLabelsForID(id string) ([]string, error) {
	if p.isDBEmpty() {
		return nil, ErrNoPasswords
	}
	var labels []string
	for key, val := range p.db.Labels {
		if utils.StringSliceContains(id, val) {
			labels = append(labels, key)
		}
	}
	return labels, nil
}

func (p *Repository) assignLabels(id string, labels []string) {
	for _, val := range labels {
		if p.isLabelExists(val) {
			p.db.Labels[val] = append(p.db.Labels[val], id)
		} else {
			p.db.Labels[val] = []string{id}
		}
	}
}

func (p *Repository) Remove(id string) error {
	if p.isDBEmpty() {
		return ErrNoPasswords
	}
	if ! p.isIDExists(id) {
		return ErrInvalidID(id)
	}
	delete(p.db.Entries, id)
	err := p.saveDB()
	if err != nil {
		return ErrCannotSavePasswordDB(err)
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

// InitRepo initialize the Password repository.
func InitRepo(mPassword string) error {
	conf, err := config.Configuration()
	if err != nil {
		return errors.Wrapf(err, "cannot get configuration")
	}
	exists, err := utils.IsFileExists(conf.PasswordDBFilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot initiate Password DB")
	}

	if !exists {
		_, err = os.Create(conf.PasswordDBFilePath)
		if err != nil {
			return errors.Wrapf(err, "unable to create Password DB file")
		}
	}
	f := &fileio.File{Path: conf.PasswordDBFilePath}
	data, err := f.Read()
	if err != nil {
		return errors.New("cannot read password DB file")
	}

	if isPasswordRepoAlreadyInitialized(data) {
		return errors.New("password repository is already initialized")
	}

	db := &DB{
		Entries: map[string]Entry{},
		Labels:  map[string][]string{},
	}
	passwordRepo := newRepository(db, mPassword, conf.EncryptorID, conf.PasswordDBFilePath)
	err = passwordRepo.saveDB()
	if err != nil {
		return errors.Wrapf(err, "unable save password repository")
	}
	return nil
}

func newRepository(db *DB, mPassword, encryptorID, dbFilePath string) *Repository {
	eFac := &encrypt.Factory{
		ID: encryptorID,
	}
	fSpec := &fileio.File{
		Path: dbFilePath,
	}
	return &Repository{
		mPassword: mPassword,
		encryptor: eFac.GetEncryptor(),
		db:        db,
		file:      fSpec,
	}
}

// LoadRepo initializes the Password repository.
func LoadRepo(mPassword, encryptorID, passwordDBFilePath string) (*Repository, error) {
	eFac := &encrypt.Factory{
		ID: encryptorID,
	}
	fSpec := &fileio.File{
		Path: passwordDBFilePath,
	}
	rawDb, err := loadDBFile(mPassword, eFac.GetEncryptor(), fSpec)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get Raw DB")
	}
	db, err := loadDB(rawDb)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get DB")
	}
	passwordRepo := &Repository{
		mPassword: mPassword,
		encryptor: eFac.GetEncryptor(),
		db:        db,
		file:      fSpec,
	}
	return passwordRepo, nil
}
