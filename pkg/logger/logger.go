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
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

type controlLogger struct {
	mu     sync.RWMutex
	logger *logger
	hooks  []func(zapcore.Entry) error
	done   atomic.Bool
}

type logger struct {
	*zap.Logger
}

var (
	control           controlLogger
	logLevel          = zapcore.InfoLevel
	callerSkip        = 2
	logFileHandler    atomic.Value
	stdErrFileHandler atomic.Value // 全局变量，避免被 GC 回收
	defaultService    = "_default"
)

// init mainly used to output logs before logger.Init
func init() {
	cfg := buildConfig(nil)
	control.logger = &logger{BuildLogger(cfg, control.addZapOptions(defaultService)...)}
	control.hooks = make([]func(zapcore.Entry) error, 0)

	klog.SetLogger(GetKlogLogger())
}

func Init(service string, level string) {
	if service == "" {
		panic("server should not be empty")
	}

	logLevel = parseLevel(level)
	control.updateLogger(service)
	control.scheduleUpdateLogger(service)
}

// AddLoggerHook 会将传进的参数在每一次日志输出后执行
func AddLoggerHook(fns ...func(zapcore.Entry) error) {
	control.hooks = append(control.hooks, fns...)
}

func (l *controlLogger) scheduleUpdateLogger(service string) {
	// 确保只开启一次定时更新
	if !l.done.Load() {
		l.done.Store(true)
		go func() {
			for {
				now := time.Now()
				//nolint
				next := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
				time.Sleep(time.Until(next))
				l.updateLogger(service)
			}
		}()
	}
}

func (l *controlLogger) updateLogger(service string) {
	// 避免 logger 更新时引发竞态
	l.mu.Lock()
	defer l.mu.Unlock()
	var err error
	var pwd string

	// 获取当前目录
	if pwd, err = getCurrentDirectory(); err != nil {
		panic(err)
	}

	// 设置文件输出的位置
	date := time.Now().Format("2006-01-02")
	logPath := fmt.Sprintf(constants.LogFilePathTemplate, pwd, constants.LogFilePath, date, service)
	stderrPath := fmt.Sprintf(constants.ErrorLogFilePathTemplate, pwd, constants.LogFilePath, date, service)

	// 打开文件,并设置无引用时关闭文件
	logFileHandler.Store(checkAndOpenFile(logPath))
	stdErrFileHandler.Store(checkAndOpenFile(stderrPath))

	// 让日志输出到不同的位置
	logLevelFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= logLevel
	})
	errLevelFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > logLevel
	})

	logCore := zapcore.NewCore(defaultEnc(), zapcore.Lock(logFileHandler.Load().(*os.File)), logLevelFn)    //nolint
	errCore := zapcore.NewCore(defaultEnc(), zapcore.Lock(stdErrFileHandler.Load().(*os.File)), errLevelFn) //nolint

	cfg := buildConfig(zapcore.NewTee(logCore, errCore))

	l.logger.Logger = BuildLogger(cfg, l.addZapOptions(service)...)
}

func (l *controlLogger) addZapOptions(serviceName string) []zap.Option {
	var opts []zap.Option
	if len(l.hooks) != 0 {
		opts = append(opts, zap.Hooks(l.hooks...))
	}
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.AddCallerSkip(callerSkip))
	opts = append(opts, zap.Fields(
		zap.String(constants.ServiceKey, serviceName),
		zap.String(constants.SourceKey, fmt.Sprintf("app-%s", serviceName)),
	))

	return opts
}

func (l *controlLogger) debug(msg string, fields ...zap.Field) {
	l.mu.RLock() // 锁的是 logger 的操作权限, 而不是写操作, 写操作在 zap.logger 的内部有锁.
	defer l.mu.RUnlock()
	l.logger.Debug(msg, fields...)
}

func (l *controlLogger) debugf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *controlLogger) info(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(msg, fields...)
}

func (l *controlLogger) infof(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *controlLogger) warn(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Warn(msg, fields...)
}

func (l *controlLogger) warnf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Warn(fmt.Sprintf(template, args...))
}

func (l *controlLogger) error(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Error(msg, fields...)
}

func (l *controlLogger) errorf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Error(fmt.Sprintf(template, args...))
}

func (l *controlLogger) fatal(msg string, fields ...zap.Field) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Fatal(msg, fields...)
}

func (l *controlLogger) fatalf(template string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.logger.Fatal(fmt.Sprintf(template, args...))
}

func parseLevel(level string) zapcore.Level {
	var lvl zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	case "fatal":
		lvl = zapcore.FatalLevel
	default:
		lvl = zapcore.InfoLevel
	}
	return lvl
}
