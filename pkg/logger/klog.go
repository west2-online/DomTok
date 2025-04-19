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
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"go.uber.org/zap"

	"github.com/west2-online/DomTok/pkg/constants"
)

type KlogLogger struct{}

func GetKlogLogger() *KlogLogger {
	return &KlogLogger{}
}

func (l *KlogLogger) Trace(v ...any) {
	control.debug(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Debug(v ...any) {
	control.debug(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Info(v ...any) {
	control.info(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Notice(v ...any) {
	control.info(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Warn(v ...any) {
	control.warn(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Error(v ...any) {
	control.error(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Fatal(v ...any) {
	control.fatal(fmt.Sprint(v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Tracef(format string, v ...any) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Debugf(format string, v ...any) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Infof(format string, v ...any) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Noticef(format string, v ...any) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Warnf(format string, v ...any) {
	control.warn(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Errorf(format string, v ...any) {
	control.error(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) Fatalf(format string, v ...any) {
	control.fatal(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxTracef(ctx context.Context, format string, v ...any) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxDebugf(ctx context.Context, format string, v ...any) {
	control.debug(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxInfof(ctx context.Context, format string, v ...any) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxNoticef(ctx context.Context, format string, v ...any) {
	control.info(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxWarnf(ctx context.Context, format string, v ...any) {
	control.warn(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxErrorf(ctx context.Context, format string, v ...any) {
	control.error(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) CtxFatalf(ctx context.Context, format string, v ...any) {
	control.fatal(fmt.Sprintf(format, v...), zap.String(constants.SourceKey, constants.KlogSource))
}

func (l *KlogLogger) SetLevel(klog.Level) {
}

func (l *KlogLogger) SetOutput(io.Writer) {
}
