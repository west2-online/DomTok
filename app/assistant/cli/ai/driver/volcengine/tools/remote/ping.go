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

package remote

import (
	"context"

	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

// Tips: This function should not be used in the future

// TODO: remove or keep as a example

// Ping pings the server
func Ping(ctx context.Context, _ string, server adapter.ServerCaller) (string, error) {
	resp, err := server.Ping(ctx)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

// if server.Ping has a request parameter, the function signature should be:
// func Ping(ctx context.Context, args string, server adapter.ServerCaller) (string, error) {
//    req := &model.PingRequest{}
//    json.Unmarshal([]byte(args), req)
//    resp, err := server.Ping(ctx, req)
//    if err != nil {
//        return "", err
//    }
//    return string(resp), nil
// }
