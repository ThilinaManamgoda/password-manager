/*
 *  Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 */

package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("SuccessTest", testReadFileSuccessFunc())
	t.Run("FailTest", testReadFileFailedFunc())
}

func testReadFileSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		wDir, err := os.Getwd()
		if err != nil {
			t.Error(err)
		}
		p := &Spec{
			Path: filepath.Join(wDir, "../../test/test_read_file_success"),
		}
		result, err := p.Read()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "test data", string(result))
	}
}

func testReadFileFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		p := &Spec{
			Path: "",
		}
		_, err := p.Read()
		if pathErr, ok := err.(*os.PathError); ok {
			assert.Equal(t, "open : no such file or directory", pathErr.Error())
		} else {
			t.Error(err)
		}
	}
}
