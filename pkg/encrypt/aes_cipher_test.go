/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

package encrypt

import (
	"github.com/ThilinaManamgoda/password-manager/pkg/file"
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
		p := &file.Spec{
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
