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

package service

import (
	"context"
	"fmt"

	"github.com/hertz-contrib/websocket"

	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/app/assistant/pack"
)

// Accept accepts a websocket message.
func (s _Service) Accept(conn *websocket.Conn, ctx context.Context) (err error) {
	// read the message from the websocket connection
	t, m, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	ctx = context.WithValue(ctx, CtxKeyInput, string(m))
	switch t {
	case websocket.TextMessage:
		// TODO: handle text message

		err := handleTextMessage(conn, ctx)
		if err != nil {
			return fmt.Errorf("handle text message failed: %w", err)
		}

	case websocket.BinaryMessage:
		// binary message is not supported
		return fmt.Errorf("binary message is not supported")

	default:
		// other message types are not expected to be handled
		return nil
	}

	return nil
}

// handleTextMessage handles a text message.
func handleTextMessage(conn *websocket.Conn, ctx context.Context) (err error) {
	// Login has been called before Accept
	id, _ := ctx.Value(CtxKeyID).(string)
	// Input is set in Accept
	input, _ := ctx.Value(CtxKeyInput).(string)
	dialog := model.NewDialog(id, input)
	errChan := make(chan error)

	err = conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command("dialog_open"))
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	defer func() {
		_ = conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command("dialog_close"))
	}()
	go func(d model.IDialog) {
		err := Service.ai.Call(ctx, d)
		errChan <- err
	}(dialog)

	index := 0
	for {
		select {
		case <-dialog.NotifyOnClosed():
			return nil
		case err := <-errChan:
			if err != nil {
				return err
			}
		case msg := <-dialog.NotifyOnMessage():
			if msg == "" {
				continue
			}
			data := map[string]interface{}{
				"index":   index,
				"content": msg,
			}
			index++
			err := conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Message(data))
			if err != nil {
				return fmt.Errorf("write failed: %w", err)
			}
		}
	}
}
