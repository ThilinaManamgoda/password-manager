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


// Package fileio handles the interaction with files
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
	Path string
}

// Read method reads the file
func (p *File) Read() ([]byte, error) {
	f, err := os.OpenFile(p.Path, os.O_CREATE|os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write method writes data to the file
func (p *File) Write(data []byte) error {
	err := ioutil.WriteFile(p.Path, data, 0640)
	if err != nil {
		return err
	}
	return nil
}

// IsFileExists checks whether the given file exists.
func IsFileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.IsDir() {
		return false, ErrPathIsADir
	}
	return true, nil
}
