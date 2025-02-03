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

package logger

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

func Debug(msg string, fields ...zap.Field) {
	control.debug(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	control.debugf(template, args...)
}

func Info(msg string, fields ...zap.Field) {
	control.info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	control.infof(template, args...)
}

func Warn(msg string, fields ...zap.Field) {
	control.warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	control.warnf(template, args...)
}

func Error(msg string, fields ...zap.Field) {
	control.error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	control.errorf(template, args...)
}

func Fatal(msg string, fields ...zap.Field) {
	control.fatal(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	control.fatalf(template, args...)
}

const permission = 0o755 // 用户具有读/写/执行权限，组用户和其它用户具有读写权限

// getCurrentDirectory 会返回当前运行的目录
func getCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(dir, "\\", "/"), nil
}

func checkAndOpenFile(path string) *os.File {
	var err error
	var handler *os.File
	if err = os.MkdirAll(filepath.Dir(path), permission); err != nil {
		panic(err)
	}

	handler, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, permission)
	if err != nil {
		panic(err)
	}
	runtime.SetFinalizer(handler, func(fd *os.File) {
		if err := fd.Close(); err != nil {
			Infof("close file failed %v", err)
		}
	})
	return handler
}
