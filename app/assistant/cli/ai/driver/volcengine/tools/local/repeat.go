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

package local

import (
	"context"
	"encoding/json"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

// Tips: This function should not be used in the future

type _RepeatArgs struct {
	Message string `json:"message"`
}

// Repeat repeats the message
// calls locally, ignore the context and server caller
func Repeat(_ context.Context, args string, _ adapter.ServerCaller) (string, error) {
	req := &_RepeatArgs{}
	err := json.Unmarshal([]byte(args), req)
	if err != nil {
		return "", err
	}
	return req.Message, nil
}
