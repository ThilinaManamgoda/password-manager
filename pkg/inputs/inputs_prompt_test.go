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

package inputs

import (
	"gotest.tools/assert"
	"testing"
)

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
