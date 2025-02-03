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

// DO NOT EDIT

package errno

import (
	"errors"
	"fmt"
	"io"
)

type ErrNo struct {
	ErrorCode int64
	ErrorMsg  string
	stack     *stack
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

func NewErrNoWithStack(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg,
		stack:     callers(),
	}
}

func Errorf(code int64, template string, args ...interface{}) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  fmt.Sprintf(template, args...),
		stack:     callers(),
	}
}

// WithMessage will replace default msg to new
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrorMsg = msg
	return e
}

// WithError will add error msg after Message
func (e ErrNo) WithError(err error) ErrNo {
	e.ErrorMsg = e.ErrorMsg + ", " + err.Error()
	return e
}

func (e ErrNo) StackTrace() any {
	if e.stack == nil { // nil 地狱
		return nil
	}
	return e.stack
}

func (e ErrNo) Format(st fmt.State, verb rune) {
	switch verb {
	case 's':
		io.WriteString(st, e.Error())
	case 'v':
		io.WriteString(st, e.Error())
		switch {
		case st.Flag('+'):
			e.stack.Format(st, verb)
		}
	}
}

// ConvertErr convert error to ErrNo
// in Default user ServiceErrorCode
func ConvertErr(err error) ErrNo {
	if err == nil {
		return Success
	}
	errno := ErrNo{}
	if errors.As(err, &errno) {
		return errno
	}

	s := InternalServiceError
	s.ErrorMsg = err.Error()
	return s
}
