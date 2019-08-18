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

package inputs

import (
	"github.com/spf13/cobra"
	"gotest.tools/assert"
	"testing"
)

func TestIsPasswordValid(t *testing.T) {
	t.Run("SuccessTest", testIsPasswordValidSuccessFunc())
	t.Run("FailedTestEmptyPassword", testIsPasswordValidFailedEmptyFunc())
}

func testIsPasswordValidSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsPasswordValid("password")
		assert.Equal(t, true, result)
	}
}

func testIsPasswordValidFailedEmptyFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsPasswordValid("")
		assert.Equal(t, false, result)
	}
}

func TestGetFlagBoolVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagBoolValSuccessFunc())
	t.Run("FailedTest", testGetFlagBoolValFailedFunc())
}

func testGetFlagBoolValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().BoolP("param", "p", false, "Test parameter")
		result, err := GetFlagBoolVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, false, result)
	}
}

func testGetFlagBoolValFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		_, err := GetFlagBoolVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestGetFlagStringVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagStringValSuccessFunc())
	t.Run("FailedTest", testGetFlagStringValFailedFunc())
}

func testGetFlagStringValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().StringP("param", "p", "DEFAULT_VAL", "Test parameter")
		result, err := GetFlagStringVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "DEFAULT_VAL", result)
	}
}

func testGetFlagStringValFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		_, err := GetFlagStringVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestGetFlagStringArrayVal(t *testing.T) {
	t.Run("SuccessTest", testGetFlagStringArrayValSuccessFunc())
	t.Run("FailedTest", testGetFlagStringArrayValFailedFunc())
}

func testGetFlagStringArrayValSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().StringArrayP("param", "p", []string{"a", "b"}, "Test parameter")
		result, err := GetFlagStringArrayVal(cmd, "param")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "a", result[0])
		assert.Equal(t, "b", result[1])
	}
}

func testGetFlagStringArrayValFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		cmd := &cobra.Command{}
		_, err := GetFlagStringArrayVal(cmd, "param")
		if err == nil {
			t.Error("Expecting an error")
		}
	}
}

func TestIsValidSingleArg(t *testing.T) {
	t.Run("SuccessTest", testIsValidSingleArgSuccessFunc())
	t.Run("FailedTestWithEmptyArgs", testIsValidSingleArgEmptyArgFailedFunc())
	t.Run("FailedTestWithMultipleArgs", testIsValidSingleArgMultipleArgsFailedFunc())
}

func testIsValidSingleArgSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, true, IsValidSingleArg([]string{"a"}))
	}
}

func testIsValidSingleArgEmptyArgFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, false, IsValidSingleArg([]string{}))
	}
}

func testIsValidSingleArgMultipleArgsFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, false, IsValidSingleArg([]string{"a", "b"}))
	}
}

func TestIsArgValid(t *testing.T) {
	t.Run("SuccessTest", testIsArgValidSuccessFunc())
	t.Run("FailedTest", testIsArgValidFailedFunc())
}

func testIsArgValidSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, true, IsArgValid("test"))
	}
}

func testIsArgValidFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, false, IsArgValid(""))
	}
}

