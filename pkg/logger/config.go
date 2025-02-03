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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	core zapcore.Core
	enc  zapcore.Encoder
	ws   zapcore.WriteSyncer
	lvl  zapcore.Level
}

func buildConfig(core zapcore.Core) *config {
	cfg := defaultConfig()
	cfg.core = core
	if cfg.core == nil {
		cfg.core = zapcore.NewCore(cfg.enc, cfg.ws, cfg.lvl)
	}

	return cfg
}

func BuildLogger(cfg *config, opts ...zap.Option) *zap.Logger {
	return zap.New(cfg.core, opts...)
}

func defaultConfig() *config {
	return &config{
		enc: defaultEnc(),
		ws:  defaultWs(),
		lvl: defaultLvl(),
	}
}

func defaultEnc() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 日志等级大写
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zapcore.NewJSONEncoder(cfg)
}

func defaultWs() zapcore.WriteSyncer {
	return os.Stdout
}

func defaultLvl() zapcore.Level {
	return zapcore.DebugLevel
}
