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

// Package fileio handles the interaction with files.
package fileio

import (
	"errors"
	"io/ioutil"
	"os"
)

// ErrPathIsADir represents an error.
var ErrPathIsADir = errors.New("path is a directory")

// File struct represent a file
type File struct {
	Path       string
	Permission os.FileMode
}

// Read method reads the file.
func (p *File) Read() ([]byte, error) {
	f, err := os.OpenFile(p.Path, os.O_CREATE|os.O_RDONLY, p.Permission)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write method writes data to the file.
func (p *File) Write(data []byte) error {
	err := ioutil.WriteFile(p.Path, data, p.Permission)
	if err != nil {
		return err
	}
	return nil
}

// IsFileExists checks whether the given file exists.
func IsFileExists(filePath string) (bool, error) {
	exists, info, err := isExits(filePath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	if info.IsDir() {
		return false, ErrPathIsADir
	}
	return true, nil
}

func isExits(path string) (bool, os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, info, nil
}

// IsDirExists checks whether the given directory exists.
func IsDirExists(dirPath string) (bool, error) {
	exists, info, err := isExits(dirPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	if !info.IsDir() {
		return false, nil
	}
	return true, nil
}

// CreateDirectory creates the given directory.
func CreateDirectory(path string) error {
	exists, err := IsDirExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
