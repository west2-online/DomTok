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

package adapter

import (
	"context"

	category "github.com/west2-online/DomTok/app/assistant/cli/ai/driver/eino/model"
	"github.com/west2-online/DomTok/app/assistant/model"
)

// AIClient is the interface for calling the AI
// It is used by the service to call the AI
type AIClient interface {
	// Call calls the AI with the dialog
	Call(ctx context.Context, dialog model.IDialog) error
	// ForgetDialog tells the AI to forget the dialog
	// This is used when the user logs out
	ForgetDialog(dialog model.IDialog)
	// SetServerCategory sets the server category
	SetServerCategory(category category.GetServerCaller)
}
