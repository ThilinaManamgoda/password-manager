/*
 * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package storage

import (
	"bytes"
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/storage/googledrive"
	"github.com/pkg/errors"
	"time"
)

const (
	// GoogleDriveStorageID is the File storage type ID.
	GoogleDriveStorageID = "GoogleDrive"
	// ConfKeyTokenFilePath represents the token file path in the conf map.
	ConfKeyTokenFilePath = "CONF_KEY_TOKEN_FILE_PATH"
	// ConfKeyDirectory represents the password file directory in the conf map.
	ConfKeyDirectory = "CONF_KEY_DIRECTORY"
	// ConfKeyPasswordDBFile represents the password Database file name in the conf map.
	ConfKeyPasswordDBFile = "CONF_KEY_PASSWORD_DB_FILE"
)

var (
	errUnableToSearchDir = func(err error, dir string) error {
		return errors.Wrapf(err, "unable to search the directory: %s in the googledrive", dir)
	}
	errUnableToSearchFile = func(err error, f string) error {
		return errors.Wrapf(err, "unable to search the password db file: %s in the googledrive", f)
	}
	errUnableToGetPasswordDBFileID = func(err error) error {
		return errors.Wrap(err, "unable to get the password db ID")
	}
)

// GoogleDrive represent a GoogleDrive as a storage.
type GoogleDrive struct {
	directory      string
	passwordDBFile string
	tokenFile      string
	client         *googledrive.Client
}

// InitForFirstTime initialises the GoogleDrive storage for the first time.
func (g *GoogleDrive) InitForFirstTime(data []byte, conf map[string]string) error {
	err := g.Init(conf)
	if err != nil {
		return err
	}

	ok, dirID, err := g.client.IsDirExists(g.directory)
	if err != nil {
		return errUnableToSearchDir(err, g.directory)
	}
	if ok {
		fileExists, _, err := g.client.IsFileExists(dirID, g.passwordDBFile)
		if err != nil {
			return errUnableToSearchFile(err, g.passwordDBFile)
		}
		if !fileExists {
			// Ignoring the created file object since no further operations are performed on it.
			_, err = g.client.CreateFile(g.passwordDBFile, googledrive.FileMimeType, bytes.NewReader(data), dirID)
			if err != nil {
				return errors.Wrapf(err, "unable to create password DB file: %s", g.passwordDBFile)
			}
		}
	} else {
		dir, err := g.client.CreateDir(g.directory, "root")
		if err != nil {
			return errors.Wrapf(err, "unable to create directory: %s", g.directory)
		}
		// Ignoring the created file object since no further operations are performed on it.
		_, err = g.client.CreateFile(g.passwordDBFile, googledrive.FileMimeType, bytes.NewReader(data), dir.Id)
		if err != nil {
			return errors.Wrapf(err, "unable to create password DB file: %s", g.passwordDBFile)
		}
	}
	return nil
}

func (g *GoogleDrive) setValues(conf map[string]string) error {
	g.directory = conf[ConfKeyDirectory]
	g.passwordDBFile = conf[ConfKeyPasswordDBFile]
	g.tokenFile = conf[ConfKeyTokenFilePath]
	if !isValidPath(g.tokenFile) {
		return errors.New("invalid path. Path cannot be empty")
	}
	return nil
}

// Init initialises the storage.
// After "InitForFirstTime" method, this function can be called to initialize the GoogleDrive storage in the consequent uses.
func (g *GoogleDrive) Init(conf map[string]string) error {
	err := g.setValues(conf)
	if err != nil {
		return ErrUnableToConfigure(err)
	}
	g.client = &googledrive.Client{
		TokenFile: g.tokenFile,
	}
	err = g.client.Init()
	if err != nil {
		return ErrUnableToConfigure(err)
	}
	return nil
}

func (g *GoogleDrive) getPasswordDBFileDirID() (string, string, error) {
	ok, dirID, err := g.client.IsDirExists(g.directory)
	if err != nil {
		return "", "", errUnableToSearchDir(err, g.directory)
	}
	if !ok {
		return "", "", errors.New(fmt.Sprintf("directory: %s doesn't exists", g.directory))
	}

	fileExists, fileID, err := g.client.IsFileExists(dirID, g.passwordDBFile)
	if err != nil {
		return "", "", errUnableToSearchFile(err, g.passwordDBFile)
	}

	if !fileExists {
		return "", "", errors.New(fmt.Sprintf("password DB file: %s doesn't exists", g.passwordDBFile))
	}
	return fileID, dirID, nil
}

// Load loads the GoogleDrive storage as a byte array.
func (g *GoogleDrive) Load() ([]byte, error) {
	fileID, _, err := g.getPasswordDBFileDirID()
	if err != nil {
		return nil, errUnableToGetPasswordDBFileID(err)
	}
	body, err := g.client.FileContent(fileID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to download password file")
	}
	return body, nil
}

// Store store data in the GoogleDrive storage.
func (g *GoogleDrive) Store(data []byte) error {
	fileID, _, err := g.getPasswordDBFileDirID()
	if err != nil {
		return errUnableToGetPasswordDBFileID(err)
	}
	return g.client.UpdateFileContent(fileID, data)
}

// Backup backups the given password database data in the GoogleDrive storage.
func (g *GoogleDrive) Backup() error {
	ok, dirID, err := g.client.IsDirExists(g.directory)
	if err != nil {
		return errUnableToSearchDir(err, g.directory)
	}
	if ok {
		fileExists, fileID, err := g.client.IsFileExists(dirID, g.passwordDBFile)
		if err != nil {
			return errUnableToSearchFile(err, g.passwordDBFile)
		}
		if !fileExists {
			return errors.Wrapf(err, "password DB file: %s doesn't exists", g.passwordDBFile)
		}
		// Ignoring the copied file object since no further operations are performed on it.
		_, err = g.client.CopyFile(fileID, g.passwordDBFile+"_backup_"+time.Now().Format("2006-01-02"),
			googledrive.FileMimeType, dirID)
	} else {
		return errors.New(fmt.Sprintf("password directory: %s doesn't exists", g.directory))
	}
	return nil
}
