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

package utils

import (
	"encoding/json"
	"github.com/sethvargo/go-password/password"
	"sort"
)

// AESEncryptID is the unique identifier for this encryptor.
const AESEncryptID = "AES"

// IsValidByteSlice method check whether the Slice is valid or not.
func IsValidByteSlice(data []byte) bool {
	return (data != nil) && (len(data) != 0)
}

// StringSliceContains check whether the specified key is in the sorted String slice.
func StringSliceContains(s []string, key string) bool {
	i := sort.SearchStrings(s, key)
	return i != len(s) && s[i] == key
}

// GeneratePassword generates a password of given length.
func GeneratePassword(len int) (string, error) {
	pass, err := password.Generate(len, len/4, len/4, false, false)
	if err != nil {
		return "", err
	}
	return pass, nil
}

// MarshalData marshals the given struct to a byte array.
func MarshalData(data interface{}) ([]byte, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return marshaledData, nil
}

// RemoveKeyFromSortedSlice removes the given key from the sorted slice.
func RemoveKeyFromSortedSlice(s []string, key string)[]string {
	i := sort.SearchStrings(s, key)
	if i == len(s) {
		return s
	}
	copy(s[i:], s[i+1:]) // Shift a[i+1:] left one index
	s[len(s)-1] = ""     // Erase last element (write zero value)
	return s[:len(s)-1]
}
