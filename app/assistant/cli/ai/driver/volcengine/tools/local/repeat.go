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
