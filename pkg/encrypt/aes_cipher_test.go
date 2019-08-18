// Copyright © 2019 Thilina Manamgoda
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this fileio except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package encrypt

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/fileio"
	"gotest.tools/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateHash(t *testing.T) {
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", createHash(""))
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", createHash("test"))
}

func TestEncrypt(t *testing.T) {
	t.Run("FailedTestInvalidPass", testEncryptFailedInvalidPassFunc())
	t.Run("FailedTestInvalidContent", testEncryptFailedInvalidContentFunc())
}

func testEncryptFailedInvalidPassFunc() func(t *testing.T) {
	return func(t *testing.T) {
		a := &AESEncryptor{}
		_, err := a.Encrypt([]byte("test"), "")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func testEncryptFailedInvalidContentFunc() func(t *testing.T) {
	return func(t *testing.T) {
		a := &AESEncryptor{}
		_, err := a.Encrypt([]byte(""), "test")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestDecrypt(t *testing.T) {
	t.Run("SuccessPass", testDecryptSuccessFunc())
	t.Run("FailedTestInvalidPass", testDecryptFailedInvalidPassFunc())
	t.Run("FailedTestInvalidContent", testDecryptFailedInvalidContentFunc())
}

func testDecryptSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		wDir, err := os.Getwd()
		if err != nil {
			t.Error(err)
		}
		p := &fileio.File{
			Path: filepath.Join(wDir, "../../test/test_decrypt_success"),
		}
		a := &AESEncryptor{}
		content, err := p.Read()
		if err != nil {
			t.Error(err)
		}
		result, err := a.Decrypt(content, "test")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "test", string(result))
	}
}

func testDecryptFailedInvalidPassFunc() func(t *testing.T) {
	return func(t *testing.T) {
		a := &AESEncryptor{}
		_, err := a.Decrypt([]byte("test"), "")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func testDecryptFailedInvalidContentFunc() func(t *testing.T) {
	return func(t *testing.T) {
		a := &AESEncryptor{}
		_, err := a.Decrypt([]byte(""), "test")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}
