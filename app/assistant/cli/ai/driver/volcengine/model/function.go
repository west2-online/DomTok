package model

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/west2-online/DomTok/app/assistant/cli/server/adapter"
)

type Function struct {
	Name        string
	Description string
	Parameters  RootParameter
	Call        func(ctx context.Context, args string, server adapter.ServerCaller) (string, error)
}

// AsTool converts the function to a tool
func (f Function) AsTool() *model.Tool {
	return &model.Tool{
		Type: model.ToolTypeFunction,
		Function: &model.FunctionDefinition{
			Name:        f.Name,
			Description: f.Description,
			Parameters:  f.Parameters,
		},
	}
}
