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
	"errors"
	"fmt"
	"sync"

	"github.com/hertz-contrib/websocket"

	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/app/assistant/pack"
	"github.com/west2-online/DomTok/pkg/errno"
)

var busy = sync.Map{}

// Accept accepts a websocket message.
func (s Core) Accept(conn *websocket.Conn, ctx context.Context) (err error) {
	// read the message from the websocket connection
	t, m, err := conn.ReadMessage()
	if err != nil {
		if errors.Is(err, websocket.ErrReadLimit) {
			return errno.NewErrNoWithStack(errno.ServiceUserCloseWebsocketConn, err.Error())
		}
		return errno.NewErrNoWithStack(errno.InternalNetworkErrorCode, err.Error())
	}

	id, _ := ctx.Value(CtxKeyID).(string)

	b, ok := busy.LoadOrStore(id, false)
	if ok {
		if isBusy, _ := b.(bool); isBusy {
			_ = conn.WriteMessage(websocket.TextMessage,
				pack.ResponseFactory.Error(fmt.Errorf("only one dialog is allowed at the same time")))
			return nil
		}
	}
	busy.Store(id, true)

	go func() {
		// set the message in the context
		ctx := context.WithValue(ctx, CtxKeyInput, string(m))
		switch t {
		case websocket.TextMessage:
			// handle text message
			err := handleTextMessage(conn, ctx)
			if err != nil {
				_ = conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Error(err))
			}

		case websocket.BinaryMessage:
			// binary message is not supported
			_ = conn.WriteMessage(websocket.TextMessage,
				pack.ResponseFactory.Error(fmt.Errorf("binary message is not supported")))

		default:
			// other message types are not expected to be handled
			_ = conn.WriteMessage(websocket.TextMessage,
				pack.ResponseFactory.Error(fmt.Errorf("unsupported message type")))
		}

		busy.Delete(id)
	}()

	return nil
}

// handleTextMessage handles a text message.
func handleTextMessage(conn *websocket.Conn, ctx context.Context) (err error) {
	// load the context
	id, input, turn := loadCtx(ctx)
	// Create a new dialog
	dialog := model.NewDialog(id, input)
	errChan := make(chan error)
	// Mark the dialog as opened
	// Therefore, the frontend can start a dialog to present the messages
	err = writeOpenDialog(conn, turn)
	if err != nil {
		if errors.Is(err, websocket.ErrCloseSent) {
			return errno.NewErrNoWithStack(errno.ServiceUserCloseWebsocketConn, err.Error())
		}
		return errno.NewErrNoWithStack(errno.InternalNetworkErrorCode, err.Error())
	}
	defer func() {
		// Mark the dialog as closed
		// Therefore, the frontend can end the dialog
		_ = writeCloseDialog(conn, turn)
	}()
	// Call the AI client asynchronously
	callAIClientWithErrChanAsync(ctx, dialog, &errChan)
	index := int64(0)
	for {
		select {
		case <-dialog.NotifyOnClosed():
			// if the dialog is closed, return nil
			return nil
		case err := <-errChan:
			// if there is an error, return it
			if err != nil {
				return err
			}
		case msg := <-dialog.NotifyOnMessage():
			// if there is a message, send it to the frontend

			// skip empty messages
			if msg == "" {
				continue
			}
			// format the message
			data := model.NewDeltaContent(msg, index, turn)
			index++
			// send the message
			err := writeMessage(conn, data)
			if err != nil {
				// if the message cannot be sent, consume the rest messages and return the error
				// to avoid goroutine leak
				go consumeRestMessages(*dialog)
				// if the error is ErrCloseSent, perhaps the connection is closed by the user
				if errors.Is(err, websocket.ErrCloseSent) {
					return errno.NewErrNoWithStack(errno.ServiceUserCloseWebsocketConn, err.Error())
				}
				return errno.NewErrNoWithStack(errno.InternalNetworkErrorCode, err.Error())
			}
		}
	}
}

// loadCtx loads the context.
func loadCtx(ctx context.Context) (string, string, int64) {
	// Login has been called before Accept
	id, _ := ctx.Value(CtxKeyID).(string)
	// Input is set in Accept
	input, _ := ctx.Value(CtxKeyInput).(string)
	// Turn is set in Service, and it is increased by 1 each time Accept is called
	// If it is not set, it is 0
	turn, ok := ctx.Value(CtxKeyTurn).(int64)
	if !ok {
		turn = 0
	}
	return id, input, turn
}

// consumeRestMessages consumes the rest messages in the dialog.
func consumeRestMessages(dialog model.Dialog) {
	for {
		select {
		case <-dialog.NotifyOnClosed():
			return
		case <-dialog.NotifyOnMessage():
		}
	}
}

// writeOpenDialog writes an open dialog message to the connection.
func writeOpenDialog(conn *websocket.Conn, turn int64) error {
	return conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command(model.NewDialogOp(pack.OpContentOpen, turn)))
}

// writeCloseDialog writes a close dialog message to the connection.
func writeCloseDialog(conn *websocket.Conn, turn int64) error {
	return conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command(model.NewDialogOp(pack.OpContentClose, turn)))
}

// writeMessage writes a message to the connection.
func writeMessage(conn *websocket.Conn, data model.DeltaContent) error {
	return conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Message(data))
}

// callAIClientWithErrChanAsync calls the AI client asynchronously.
func callAIClientWithErrChanAsync(ctx context.Context, d model.IDialog, errChan *chan error) {
	go func() {
		err := Service.ai.Call(ctx, d)
		*errChan <- err
	}()
}
