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

package storage

import (
	"fmt"
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	// ConfKeyFilePath represents the password db path key in the conf map.
	ConfKeyFilePath = "CONF_KEY_FILE_PATH"
	// ConfKeyFilePermission represents the password db file permission key in the conf map.
	ConfKeyFilePermission = "CONF_KEY_FILE_PERMISSION"
	// FilePermissionPattern represents the password db file permission regex.
	FilePermissionPattern = "[0,2,4,6][0-7][0-7][0-7]"
	// FileStorageID is the File storage type ID.
	FileStorageID = "File"
)

var (
	// ErrInvalidPermission represents an invalid permission error.
	ErrInvalidPermission = func(p string) error { return errors.New(fmt.Sprintf("invalid permission: %s", p)) }
	// ErrUnableToConfigure represents unable to configure file storage error.
	ErrUnableToConfigure = func(e error) error { return errors.Wrap(e, "unable to configure storage") }
)

// File represent a file as a storage.
type File struct {
	path       string
	permission os.FileMode
}

// InitForFirstTime initialises the file storage for the first time.
func (f *File) InitForFirstTime(data []byte, conf map[string]string) error {
	err := f.setValues(conf)
	if err != nil {
		return ErrUnableToConfigure(err)
	}
	exists, err := fileio.IsFileExists(f.path)
	if err != nil {
		return errors.Wrapf(err, "cannot inspect the file storage")
	}
	if exists {
		return errors.New(fmt.Sprintf("file: %s already exists", f.path))
	}

	// Ignoring the returned pointer to os.File since it is not required.
	_, err = os.Create(f.path)
	if err != nil {
		return errors.Wrapf(err, "unable to create storage file")
	}

	err = f.Store(data)
	if err != nil {
		return err
	}
	return nil
}

// Init initialises the file storage.
func (f *File) Init(conf map[string]string) error {
	err := f.setValues(conf)
	if err != nil {
		return ErrUnableToConfigure(err)
	}
	return nil
}

// Load loads the file storage as a byte array.
func (f *File) Load() ([]byte, error) {
	file := fileio.File{Path: f.path, Permission: f.permission}
	data, err := file.Read()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load the file")
	}
	return data, nil
}

// Store store data in the file storage.
func (f *File) Store(data []byte) error {
	err := store(f.path, f.permission, data)
	if err != nil {
		return errors.Wrap(err, "unable store file")
	}
	return nil
}

// Backup backups file in the file storage.
func (f *File) Backup() error {
	data, err := f.Load()
	if err != nil {
		return errors.Wrap(err, "unable read current file")
	}
	backupFilePath := f.path + "_backup_" + time.Now().Format("2006-01-02")
	exists, err := fileio.IsFileExists(backupFilePath)
	if err != nil {
		return errors.Wrapf(err, "cannot inspect the file storage for backup")
	}
	if exists {
		return errors.New(fmt.Sprintf("backup file: %s already exists", backupFilePath))
	}
	// Ignoring the returned pointer to os.File since it is not required.
	_, err = os.Create(backupFilePath)
	if err != nil {
		return errors.Wrapf(err, "unable to create storage file")
	}
	err = store(backupFilePath, f.permission, data)
	if err != nil {
		return errors.Wrap(err, "unable backup file")
	}
	return nil
}

func store(path string, permission os.FileMode, data []byte) error {
	file := fileio.File{Path: path, Permission: permission}
	err := file.Write(data)
	if err != nil {
		return errors.Wrap(err, "unable to write data to file")
	}
	return nil
}

func (f *File) setValues(conf map[string]string) error {
	path := conf[ConfKeyFilePath]
	if !isValidPath(path) {
		return errors.New("invalid path. Path cannot be empty")
	}
	f.path = path

	permission := conf[ConfKeyFilePermission]

	p, err := getPermission(permission)
	if err != nil {
		return err
	}
	f.permission = p
	return nil
}

func isValidPath(p string) bool {
	return p != ""
}

func getPermission(p string) (os.FileMode, error) {
	if len(p) != 4 {
		return 0, ErrInvalidPermission(p)
	}
	match, err := regexp.MatchString(FilePermissionPattern, p)
	if err != nil {
		return 0, errors.Wrap(err, "unable to validate permission against the pattern")
	}
	if !match {
		return 0, ErrInvalidPermission(p)
	}
	u, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "unable to convert permission to uint type")
	}
	return os.FileMode(u), nil
}
