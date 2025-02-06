package template

import (
	"context"

	"github.com/west2-online/DomTok/app/user/entities"
)

type TemplateRPCClient struct {
	// client templateservice.client
}

func NewTemplateRPCClient() *TemplateRPCClient {
	return &TemplateRPCClient{}
}

func (c *TemplateRPCClient) GetTemplateInfo(ctx context.Context) (*entities.TemplateModel, error) {
	// rpc call
	// client.Call()
	return nil, nil
}
