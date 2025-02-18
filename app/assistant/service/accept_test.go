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
	"sync"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hertz-contrib/websocket"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/model"
	"github.com/west2-online/DomTok/app/assistant/pack"
)

func TestCore_Accept(t *testing.T) {
	ctx := context.Background()
	type Message struct {
		t int
		m []byte
	}
	msgChan := make(chan Message, 1)
	delChan := make(chan string, 1)
	successErr := errors.New("success")
	Normalize := func() {
		ctx = context.WithValue(ctx, CtxKeyID, "test")
		Mock((*websocket.Conn).ReadMessage).Return(websocket.TextMessage, []byte{}, nil).Build()
		Mock((*websocket.Conn).WriteMessage).To(func(t int, msg []byte) error {
			msgChan <- Message{t, msg}
			return nil
		}).Build()
		Mock((*sync.Map).Delete).To(func(key interface{}) {
			delChan <- key.(string) //nolint:forcetypeassert
		}).Build()
		Mock(handleTextMessage).Return(successErr).Build()
	}

	PatchConvey("Test Core.Accept", t, func() {
		Normalize()
		PatchConvey("success", func() {
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldBeNil)
			msg := <-msgChan
			So(msg.t, ShouldEqual, websocket.TextMessage)
			So(msg.m, ShouldEqual, pack.ResponseFactory.Error(successErr))
			key := <-delChan
			So(key, ShouldEqual, "test")
		})

		PatchConvey("read err: websocket.ErrReadLimit", func() {
			busy.Clear()
			Mock((*websocket.Conn).ReadMessage).Return(websocket.TextMessage, []byte{}, websocket.ErrReadLimit).Build()
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldNotBeNil)
			_, exist := busy.Load("test")
			So(exist, ShouldBeFalse)
		})

		PatchConvey("read err: other err", func() {
			busy.Clear()
			Mock((*websocket.Conn).ReadMessage).Return(websocket.TextMessage, []byte{}, errors.New("read err")).Build()
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldNotBeNil)
			_, exist := busy.Load("test")
			So(exist, ShouldBeFalse)
		})

		PatchConvey("busy", func() {
			busy.Clear()
			busy.Store("test", true)
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldBeNil)
			msg := <-msgChan
			So(msg.t, ShouldEqual, websocket.TextMessage)
			So(len(msg.m), ShouldBeGreaterThan, 0)
		})

		PatchConvey("binary message", func() {
			busy.Clear()
			Mock((*websocket.Conn).ReadMessage).Return(websocket.BinaryMessage, []byte{}, nil).Build()
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldBeNil)
			msg := <-msgChan
			So(msg.t, ShouldEqual, websocket.TextMessage)
			So(len(msg.m), ShouldBeGreaterThan, 0)
			key := <-delChan
			So(key, ShouldEqual, "test")
		})

		PatchConvey("unsupported message type", func() {
			busy.Clear()
			Mock((*websocket.Conn).ReadMessage).Return(3, []byte{}, nil).Build()
			err := Core{}.Accept(nil, ctx)
			So(err, ShouldBeNil)
			msg := <-msgChan
			So(msg.t, ShouldEqual, websocket.TextMessage)
			So(len(msg.m), ShouldBeGreaterThan, 0)
			key := <-delChan
			So(key, ShouldEqual, "test")
		})
	})
}

func TestHandleMessage(t *testing.T) {
	ctx := context.Background()
	Normalize := func() {
		ctx = context.WithValue(ctx, CtxKeyID, "test")
		ctx = context.WithValue(ctx, CtxKeyInput, "input")
		ctx = context.WithValue(ctx, CtxKeyTurn, 1)
		Mock((*websocket.Conn).WriteMessage).Return(nil).Build()
		Mock(writeOpenDialog).Return(nil).Build()
		Mock(writeCloseDialog).Return(nil).Build()
		Mock(writeMessage).Return(nil).Build()
		Mock(consumeRestMessages).Return().Build()
		Mock(callAIClientWithErrChanAsync).Return().Build()
	}

	PatchConvey("Test handleTextMessage", t, func() {
		Normalize()
		PatchConvey("success", func() {
			d := model.NewDialog("test", "input")
			Mock(model.NewDialog).Return(d).Build()
			d.Close()
			err := handleTextMessage(nil, ctx)
			So(err, ShouldBeNil)
		})

		PatchConvey("write open dialog err", func() {
			Mock(writeOpenDialog).Return(errors.New("write open dialog err")).Build()
			err := handleTextMessage(nil, ctx)
			So(err, ShouldNotBeNil)
		})

		PatchConvey("call AI client err", func() {
			Mock(callAIClientWithErrChanAsync).To(func(_ context.Context, d model.IDialog, ch *chan error) {
				go func() {
					*ch <- errors.New("call AI client err")
				}()
			}).Build()
			d := model.NewDialog("test", "input")
			Mock(model.NewDialog).Return(d).Build()
			err := handleTextMessage(nil, ctx)
			So(err, ShouldNotBeNil)
		})

		PatchConvey("write msg err", func() {
			Mock(writeMessage).Return(errors.New("write msg err")).Build()
			d := model.NewDialog("test", "input")
			Mock(model.NewDialog).Return(d).Build()
			d.Send("msg")
			err := handleTextMessage(nil, ctx)
			So(err, ShouldNotBeNil)
		})
	})
}
