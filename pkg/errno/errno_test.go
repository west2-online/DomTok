/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errno

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	code   = 200
	msg    = "ok"
	sucMsg = "[200] ok"
)

func getErr() ErrNo {
	return NewErrNo(code, msg)
}

func TestNewErrno(t *testing.T) {
	err := getErr()
	assert.Equal(t, sucMsg, err.Error())
	assert.Nil(t, err.StackTrace())
}

func TestNewErrnoWithStack(t *testing.T) {
	err := NewErrNoWithStack(code, msg)
	assert.Equal(t, sucMsg, err.Error())
	assert.NotNil(t, err.StackTrace())
}

func TestErrorf(t *testing.T) {
	err := Errorf(code, msg)
	assert.Equal(t, sucMsg, err.Error())
	assert.NotNil(t, err.StackTrace())
}

func TestWithMessage(t *testing.T) {
	err := getErr()
	err = err.WithMessage("success")
	assert.Equal(t, "[200] success", err.Error())
}

func TestWithError(t *testing.T) {
	err := getErr()
	err = err.WithError(fmt.Errorf("success"))
	assert.Equal(t, "[200] ok, success", err.Error())
}

func TestConvertErrWhenNil(t *testing.T) {
	suc := ConvertErr(nil)
	assert.Equal(t, Success, suc)
}

func TestConvertErrWhenErrNo(t *testing.T) {
	originErrno := getErr()
	packError := fmt.Errorf("pack, %w", originErrno)
	origin := ConvertErr(packError)
	assert.Equal(t, originErrno, origin)
}

func TestConvertErrWhenNormal(t *testing.T) {
	normal := fmt.Errorf("normal")
	err := ConvertErr(normal)
	assert.Equal(t, InternalServiceError.ErrorCode, err.ErrorCode)
	assert.Equal(t, "normal", err.ErrorMsg)
}

func TestFormat(t *testing.T) {
	err := getErr()
	assert.Equal(t, "", fmt.Sprintf("%d", err))
	assert.NotEqual(t, "", fmt.Sprintf("%s", err))
	assert.NotEqual(t, "", fmt.Sprintf("%v", err))
	assert.NotEqual(t, "", fmt.Sprintf("%+v", err))
}
