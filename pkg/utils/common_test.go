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
	"github.com/spf13/cobra"
	"gotest.tools/assert"
	"testing"
)

func TestIsValidByteSlice(t *testing.T) {
	t.Run("SuccessTest", testIsValidByteSliceSuccessFunc())
	t.Run("FailTestEmptySlice", testIsValidByteSliceFailEmptySliceFunc())
	t.Run("FailTestNilSlice", testIsValidByteSliceFailNilSliceFunc())
}

func testIsValidByteSliceSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice([]byte("test"))
		assert.Equal(t,true ,result)
	}
}

func testIsValidByteSliceFailEmptySliceFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice([]byte(""))
		assert.Equal(t,false ,result)
	}
}

func testIsValidByteSliceFailNilSliceFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice(nil)
		assert.Equal(t,false ,result)
	}
}

func TestIsPasswordValid(t *testing.T) {
	t.Run("SuccessTest", testIsPasswordValidSuccessFunc())
	t.Run("FailTestEmptyPassword", testIsValidByteSliceFailEmptySliceFunc())
	t.Run("FailTestEmptyPassword", testIsPasswordValidFailsFunc())
}

func testIsPasswordValidSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsPasswordValid("password")
		assert.Equal(t,true ,result)
	}
}

func testIsPasswordValidFailsFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsPasswordValid("")
		assert.Equal(t,false ,result)
	}
}

func TestGetFlagBoolVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagBoolValSuccessFunc())
	t.Run("FailTest",testGetFlagBoolValFailsFunc() )
	//t.Run("FailTestEmptyPassword", )
}

func testGetFlagBoolValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		cmd.Flags().BoolP("param", "p", false,"Test parameter")
		result, err := GetFlagBoolVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, false, result)
	}
}

func testGetFlagBoolValFailsFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		_, err := GetFlagBoolVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestGetFlagStringVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagStringValSuccessFunc())
	t.Run("FailTest",testGetFlagStringValFailsFunc() )
	//t.Run("FailTestEmptyPassword", )
}

func testGetFlagStringValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		cmd.Flags().StringP("param", "p", "DEFAULT_VAL","Test parameter")
		result, err := GetFlagStringVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "DEFAULT_VAL", result)
	}
}

func testGetFlagStringValFailsFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		_, err := GetFlagStringVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestGetFlagStringArrayVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagStringArrayValSuccessFunc())
	t.Run("FailTest",testGetFlagStringArrayValFailsFunc() )
}

func testGetFlagStringArrayValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		cmd.Flags().StringArrayP("param", "p", []string{"a", "b"},"Test parameter")
		result, err := GetFlagStringArrayVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "a", result[0])
		assert.Equal(t, "b", result[1])
	}
}

func testGetFlagStringArrayValFailsFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{
		}
		_, err := GetFlagStringArrayVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestConfiguration(t *testing.T) {
	defaultConf, err := defaultConf()
	if err != nil {
		t.Error(err)
	}
	result, err := Configuration()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t,defaultConf.PasswordFilePath, result.PasswordFilePath)
	assert.Equal(t,defaultConf.EncryptorID, result.EncryptorID)
}