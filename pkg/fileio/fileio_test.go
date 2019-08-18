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


package fileio

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
		p := &File{
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
		p := &File{
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
