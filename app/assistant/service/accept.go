package service

import (
	"fmt"
	"github.com/hertz-contrib/websocket"
	"github.com/west2-online/DomTok/app/assistant/cli/ai"
	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/app/assistant/pack"
)

// Accept accepts a websocket message.
func (s Service) Accept(conn *websocket.Conn) (err error) {
	// read the message from the websocket connection
	t, m, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read failed: %v", err)
	}

	switch t {
	case websocket.TextMessage:
		// TODO: handle text message

		err := handleTextMessage(conn, string(m))
		if err != nil {
			return fmt.Errorf("handle text message failed: %v", err)
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
func handleTextMessage(conn *websocket.Conn, input string) (err error) {
	dialog := model.NewDialog()
	errChan := make(chan error)

	err = conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command("dialog_open"))
	if err != nil {
		return fmt.Errorf("write failed: %v", err)
	}
	defer func() {
		_ = conn.WriteMessage(websocket.TextMessage, pack.ResponseFactory.Command("dialog_close"))
	}()
	go func(d model.IDialog) {
		errChan <- ai.Example(input, d)
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
				return fmt.Errorf("write failed: %v", err)
			}
		}
	}
}
