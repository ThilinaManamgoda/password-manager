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

// Package passwords contains the in memory data manipulations for password repo.
package passwords

import (
	"encoding/json"
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/config"
	"github.com/ThilinaManamgoda/password-manager/pkg/encrypt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	storage2 "github.com/ThilinaManamgoda/password-manager/pkg/storage"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"strings"
)

// Entry struct represents entry in the password db.
type Entry struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description,omitempty"`
}

var (
	// ErrInvalidID represents the invalid ID error.
	ErrInvalidID = func(id string) error {
		return errors.New(fmt.Sprintf("Invalid ID:  %s", id))
	}
	// ErrCannotSavePasswordDB represents the error when it is unable to save password entry.
	ErrCannotSavePasswordDB = func(err error) error {
		return errors.Wrap(err, "cannot save password db")
	}
	// ErrNoPasswords represents the error when no passwords are available.
	ErrNoPasswords = errors.New("no passwords are available")
	// ErrCannotFindMatchForID represents the error when it is unable to find a password entry for the given ID.
	ErrCannotFindMatchForID = func(id string) error {
		return errors.New(fmt.Sprintf("cannot find any match for id %s", id))
	}
)

// DB struct represents password db.
type DB struct {
	Entries map[string]Entry    `json:"entries"`
	Labels  map[string][]string `json:"labels"`
}

// Repository struct handles Password db.
type Repository struct {
	encryptor encrypt.Encryptor
	mPassword string
	db        *DB
	storage   storage2.Storage
}

func isPasswordRepoAlreadyInitialized(repoData []byte) bool {
	return utils.IsValidByteSlice(repoData)
}

func (p *Repository) marshalDB() ([]byte, error) {
	passwordDBJSON, err := utils.MarshalData(p.db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal the password db")
	}
	return passwordDBJSON, nil
}

// ChangeMasterPassword changes the master password.
func (p *Repository) ChangeMasterPassword(newPassword string) error {
	p.mPassword = newPassword
	err := p.saveDB()
	if err != nil {
		return errors.Wrap(err, "cannot save password DB")
	}
	return nil
}

func (p *Repository) getEncryptedDB() ([]byte, error) {
	passwordDBJSON, err := p.marshalDB()
	if err != nil {
		return nil, err
	}
	encryptedData, err := p.encryptor.Encrypt(passwordDBJSON, p.mPassword)
	if err != nil {
		return nil, errors.Wrap(err, "cannot encrypt password db")
	}
	return encryptedData, nil
}

func (p *Repository) saveDB() error {
	encryptedData, err := p.getEncryptedDB()
	if err != nil {
		return errors.Wrap(err, "cannot encrypt password db")
	}
	err = p.storage.Store(encryptedData)
	if err != nil {
		return ErrCannotSavePasswordDB(err)
	}
	return nil
}

func (p *Repository) addPasswordEntryToRepo(id, uN, password, desc string, labels []string) error {
	if p.isIDExists(id) {
		return errors.New(fmt.Sprintf("ID: %s is already there !", id))
	}
	entries := p.db.Entries
	entries[id] = Entry{
		ID:          id,
		Username:    uN,
		Password:    password,
		Description: desc,
	}
	p.assignLabels(id, labels)
	return nil
}

// Add method add new password entry to Password db.
func (p *Repository) Add(id, uN, password, desc string, labels []string) error {
	if id == "" {
		return errors.New("invalid the ID")
	}
	err := p.addPasswordEntryToRepo(id, uN, password, desc, labels)
	if err != nil {
		return err
	}

	err = p.saveDB()
	if err != nil {
		return err
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

// GetUsernamePassword method retrieves username and password from Password db.
func (p *Repository) GetUsernamePassword(id string, showPassword bool) error {
	passwordEntry, err := p.GetPasswordEntry(id)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Username: %s", passwordEntry.Username))
	fmt.Println(fmt.Sprintf("Description: %s", passwordEntry.Description))
	if showPassword {
		fmt.Println(fmt.Sprintf("Password: %s", passwordEntry.Password))
	} else {
		err := clipboard.WriteAll(passwordEntry.Password)
		if err != nil {
			return errors.Wrap(err, "cannot write to clip board")
		}
		fmt.Println("Password is copied to the clip board")
	}
	return nil
}

// GetPasswordEntry gets the password entry from password db.
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

// ChangePasswordEntry changes the password entry.
func (p *Repository) ChangePasswordEntry(id string, entry Entry) error {
	passwordDB := p.db.Entries
	if len(passwordDB) == 0 {
		return ErrNoPasswords
	}
	passwordDB[id] = entry
	err := p.saveDB()
	if err != nil {
		return err
	}
	return nil
}

// SearchEntriesByID will return the password entries if the password ID contains the provide key.
func (p *Repository) SearchEntriesByID(id string) ([]Entry, error) {
	if p.isDBEmpty() {
		return nil, ErrNoPasswords
	}
	var result []Entry
	for key, val := range p.db.Entries {
		if strings.Contains(key, id) {
			result = append(result, val)
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

// SearchLabel will return the password entries if the password labels contains the provide label.
func (p *Repository) SearchLabel(label string) ([]Entry, error) {
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
	var entries []Entry
	for _, val := range uniqueIDs {
		entries = append(entries, p.db.Entries[val])
	}
	return entries, nil
}

// searchLabelsForID will return the list of labels for given ID.
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

// Remove removes the password entry of the given id.
func (p *Repository) Remove(id string) error {
	if p.isDBEmpty() {
		return ErrNoPasswords
	}
	if !p.isIDExists(id) {
		return ErrInvalidID(id)
	}
	delete(p.db.Entries, id)
	err := p.saveDB()
	if err != nil {
		return err
	}
	return nil
}

func (p *Repository) loadDB() error {
	encryptedData, err := p.storage.Load()
	if err != nil {
		return errors.Wrap(err, "unable to password DB")
	}
	if !isPasswordRepoAlreadyInitialized(encryptedData) {
		return errors.New("password repository is not initialized")
	}

	decryptedData, err := p.encryptor.Decrypt(encryptedData, p.mPassword)
	if err != nil {
		return errors.Wrap(err, "cannot decrypt password db")
	}

	var db DB
	if err := json.Unmarshal(decryptedData, &db); err != nil {
		return errors.Wrapf(err, "cannot unmarshal password repository")
	}
	p.db = &db
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

	db := &DB{
		Entries: map[string]Entry{},
		Labels:  map[string][]string{},
	}
	eFac := &encrypt.Factory{
		ID: conf.EncryptorID,
	}
	sFac := storage2.Factory{ID: conf.StorageID}
	repo := &Repository{
		mPassword: mPassword,
		encryptor: eFac.Encryptor(),
		db:        db,
		storage:   sFac.Storage(),
	}
	encryptedData, err := repo.getEncryptedDB()
	if err != nil {
		return errors.Wrap(err, "cannot encrypt password db")
	}

	err = createConfigDir(conf.DirectoryPath)
	if err != nil {
		return errors.Wrap(err, "unable to create configuration directory")
	}

	err = repo.storage.InitForFirstTime(encryptedData, conf.Storage)
	if err != nil {
		return errors.Wrapf(err, "unable initiate password repository")
	}
	return nil
}

func createConfigDir(path string) error {
	err := fileio.CreateDirectory(path)
	if err != nil {
		return err
	}
	return nil
}

// LoadRepo initializes the Password repository.
// This function initialises each component of the password repository.
func LoadRepo(mPassword string) (*Repository, error) {
	conf, err := config.Configuration()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get configuration")
	}

	eFac := &encrypt.Factory{
		ID: conf.EncryptorID,
	}
	sFac := storage2.Factory{ID: conf.StorageID}
	repo := &Repository{
		encryptor: eFac.Encryptor(),
		storage:   sFac.Storage(),
		mPassword: mPassword,
	}

	err = repo.storage.Init(conf.Storage)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to init password storage")
	}
	err = repo.loadDB()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load password DB")
	}
	return repo, nil
}
