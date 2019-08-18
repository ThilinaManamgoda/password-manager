/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

package config

import (
	"gotest.tools/assert"
	"testing"
)

func TestConfiguration(t *testing.T) {
	defaultConf, err := defaultConf()
	if err != nil {
		t.Error(err)
	}
	result, err := Configuration()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, defaultConf.PasswordFilePath, result.PasswordFilePath)
	assert.Equal(t, defaultConf.EncryptorID, result.EncryptorID)
}

