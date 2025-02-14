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

	"github.com/west2-online/DomTok/app/assistant/cli/ai/driver/volcengine/model"
	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
	"github.com/west2-online/DomTok/pkg/logger"
)

// Tips: This function should not be used in the future

// TODO: remove or keep as a example

func Ping() model.Function {
	return model.Function{
		Name:        "ping",
		Description: "测试服务是否正在运行",
		Parameters: model.RootParameter{
			Type: model.RootType,
			Properties: map[string]model.PropertyField{
				"argument": model.BaseProperty{
					Type:        model.StringType,
					Description: "参数,设定的角色名",
				},
			},
			Required: []string{"argument"},
		},
		Call: ping,
	}
}

// ping pings the server
func ping(ctx context.Context, args string, server adapter.ServerCaller) (string, error) {
	resp, err := server.Ping(ctx)
	if err != nil {
		return "", err
	}

	// TODO: remove this line
	logger.Infof(`{
Stage: "remote.Ping",
args: %v,
resp: %v,
err: %v,
}`, args, string(resp), err)
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
