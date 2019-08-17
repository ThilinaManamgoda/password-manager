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

// Package encrypt holds required functionality for encryption and decryption
package encrypt

import "github.com/password-manager/pkg/utils"

// encryptor interface exposes functions for encrypt/decrypt
type Encryptor interface {
	Encrypt(data []byte, passphrase string) ([]byte, error)
	Decrypt(data []byte, passphrase string) ([]byte, error)
}

// Factory struct holds Encrypts
type Factory struct {
	ID string
}

// GetEncryptor method returns an encryptor
func (f *Factory) GetEncryptor() Encryptor {
	switch f.ID {
	case utils.AESEncryptID:
		return &AESEncryptor{}
	default:
		return &AESEncryptor{}
	}
}
