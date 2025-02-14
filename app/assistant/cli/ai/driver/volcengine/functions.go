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

package volcengine

import (
	"context"
	"fmt"
	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/model"
	"github.com/west2-online/DomTok/pkg/logger"

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/tools/remote"
	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

var functions *[]model.Function

// GetFunctions returns the functions
func GetFunctions() *[]model.Function {
	if functions != nil {
		return functions
	}
	RebuildFunctions()
	return functions
}

// RebuildFunctions rebuilds the functions
func RebuildFunctions() {
	fs := BuildFunctions()
	functions = &fs
}

// CallFunction calls a function by name
func CallFunction(ctx context.Context, name, args string, server adapter.ServerCaller) (string, error) {
	// TODO: remove this line
	logger.Infof(`{
Stage: "CallFunction",
Name: %s,
Args: %s,
}`, name, args)
	for _, f := range *GetFunctions() {
		if f.Name == name {
			return f.Call(ctx, args, server)
		}
	}
	return "", fmt.Errorf("function %s not found", name)
}

// BuildFunctions builds the functions
func BuildFunctions() []model.Function {
	return []model.Function{
		remote.Ping(),
	}
}
