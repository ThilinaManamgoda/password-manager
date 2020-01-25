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
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"regexp"
)

const (
	ConfKeyPath           = "CONF_KEY_PATH"
	ConfKeyPermission     = "CONF_KEY_PERMISSION"
	DefaultFilePermission = 0640
	FilePermissionPattern = "[0,2,4,6][0-7][0-7][0-7]"
)

type File struct {
	path       string
	permission os.FileMode //0640
}

func (f *File) InitForFirstTime(conf map[string]string) error {
	err := f.setValues(conf)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Init(conf map[string]string) error {

	return nil
}

func (f *File) Load() ([]byte, error) {
	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_RDONLY, f.permission)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open file")
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read file")
	}
	return data, nil
}

func (f *File) Store(data []byte) error {
	err := ioutil.WriteFile(f.path, data, f.permission)
	if err != nil {
		return errors.Wrap(err, "unable to write to file")
	}
	return nil
}

func (f *File) setValues(conf map[string]string) error {
	path := conf[ConfKeyPath]
	if !isValidPath(path) {
		return errors.New("invalid path. Path cannot be empty")
	}
	f.path = path

	permission := conf[ConfKeyPermission]
	if !isPermissionConfigured(permission) {
		f.permission = DefaultFilePermission
	} else {
		f.permission = os.FileMode()
	}
}

func isValidPath(p string) bool {
	return p != ""
}

func isPermissionConfigured(p string) bool {
	return p != ""
}

func validatePermission(p string) (bool, error) {
	if len(p) == 4 {
		match, err := regexp.MatchString(FilePermissionPattern, p)
		if err != nil {
			return false, errors.Wrap(err, "unable to validate permission against the pattern")
		}
		return match, nil
	}
	return false, nil
}
