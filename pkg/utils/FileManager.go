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

package utils

import (
	"io/ioutil"
	"os"
)

type PasswordFile struct {
	File string
}

func (p *PasswordFile) ReadFile() ([]byte, error) {
	f, err := os.OpenFile(p.File, os.O_CREATE|os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PasswordFile) WriteToFile(data []byte) error {
	err := ioutil.WriteFile(p.File, data,0640 )
	if err != nil {
		return err
	}
	return nil
}
