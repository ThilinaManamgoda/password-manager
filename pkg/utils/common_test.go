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
	"gotest.tools/assert"
	"sort"
	"testing"
)

func TestIsValidByteSlice(t *testing.T) {
	t.Run("SuccessTest", testIsValidByteSliceSuccessFunc())
	t.Run("FailedTestEmptySlice", testIsValidByteSliceFailEmptySliceFunc())
	t.Run("FailedTestNilSlice", testIsValidByteSliceFailNilSliceFunc())
}

func testIsValidByteSliceSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice([]byte("test"))
		assert.Equal(t, true, result)
	}
}

func testIsValidByteSliceFailEmptySliceFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice([]byte(""))
		assert.Equal(t, false, result)
	}
}

func testIsValidByteSliceFailNilSliceFunc() func(t *testing.T) {
	return func(t *testing.T) {
		result := IsValidByteSlice(nil)
		assert.Equal(t, false, result)
	}
}

func TestStringSliceContains(t *testing.T) {
	t.Run("SuccessTest", testStringSliceContainsSuccessFunc())
	t.Run("FailedTest", testStringSliceContainsFailedFunc())
}

func testStringSliceContainsSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, true, StringSliceContains([]string{"test"}, "test"))
	}
}

func testStringSliceContainsFailedFunc() func(t *testing.T) {
	return func(t *testing.T) {
		assert.Equal(t, false, StringSliceContains([]string{"test"}, "invalid-key"))
	}
}

func TestRemoveKeyFromSortedSlice(t *testing.T) {
	arr1 := []string{"a", "b", "c"}
	sort.Strings(arr1)
	assert.Equal(t,true, sort.StringsAreSorted(arr1))
	removeB:= RemoveKeyFromSortedSlice(arr1, "b")
	assert.Equal(t, 2, len(removeB))
	removeD:= RemoveKeyFromSortedSlice(removeB, "d")
	assert.Equal(t, 2, len(removeD))
	removeC:= RemoveKeyFromSortedSlice(removeD, "c")
	assert.Equal(t, 1, len(removeC))
	assert.Equal(t, "a", removeC[0])
	removeA:= RemoveKeyFromSortedSlice(removeC, "a")
	assert.Equal(t, 0, len(removeA))
}
