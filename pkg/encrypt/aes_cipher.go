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

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ThilinaManamgoda/password-manager/pkg/utils"
	"io"
)

// AESEncryptor struct represent the data needed for AES encryption and decryption.
type AESEncryptor struct {
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt method encrypts the given data
func (a *AESEncryptor) Encrypt(data []byte, passphrase string) ([]byte, error) {
	if !utils.IsPasswordValid(passphrase) {
		return nil, errors.New("invalid password")
	}
	if !utils.IsValidByteSlice(data) {
		return nil, errors.New("invalid content")
	}
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt method decrypts the given data
func (a *AESEncryptor) Decrypt(data []byte, passphrase string) ([]byte, error) {
	if !utils.IsPasswordValid(passphrase) {
		return nil, errors.New("invalid password")
	}
	if !utils.IsValidByteSlice(data) {
		return nil, errors.New("invalid content")
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
