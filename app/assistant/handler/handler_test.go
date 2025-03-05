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

package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/service"
	"github.com/west2-online/DomTok/pkg/constants"
)

func EntryPointNormalize() {
	Mock((*websocket.HertzUpgrader).Upgrade).Return(nil).Build()
	Mock((*service.Core).Login).Return(nil).Build()
	Mock((*service.Core).Logout).Return(nil).Build()
	Mock((*service.Core).Accept).Return(Sequence(nil).Then(errors.New(""))).Build()
}

func CatchErrInRequestCtx() error {
	errChan := make(chan error)
	var err error
	Mock((*app.RequestContext).JSON).To(func(code int, obj interface{}) {
		if err, ok := obj.(error); ok {
			errChan <- err
		}
	}).Build()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		Entrypoint(context.Background(), &app.RequestContext{})
		cancel()
	}()
	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

func TestEntrypoint(t *testing.T) {
	PatchConvey("Test Entrypoint", t, func() {
		EntryPointNormalize()

		PatchConvey("on success", func() {
			err := CatchErrInRequestCtx()
			So(err, ShouldBeNil)
		})

		PatchConvey("on upgrade error", func() {
			Mock((*websocket.HertzUpgrader).Upgrade).Return(fmt.Errorf("upgrade error")).Build()
			err := CatchErrInRequestCtx()
			So(err, ShouldNotBeNil)
		})
	})
}

func WebsocketHandlerNormalize() {
	Mock((*app.RequestContext).GetHeader).To(func(key string) []byte {
		if key == constants.AuthHeader {
			return []byte("token")
		}
		return []byte("")
	}).Build()
	Mock(writePingMessage).Return(nil).Build()
	Mock(service.Core.Accept).Return(io.EOF).Build()
	Mock(service.Core.Login).Return(nil).Build()
	Mock(service.Core.Logout).Return(nil).Build()
}

func CatchErrInWebsocketConn(ctx context.Context, c *app.RequestContext) error {
	errChan := make(chan error)
	var err error
	Mock(writeError).To(func(conn *websocket.Conn, e error) {
		errChan <- e
	}).Build()
	Mock(writeTokenError).To(func(conn *websocket.Conn) {
		errChan <- errors.New("token error")
	}).Build()
	done, cancel := context.WithCancel(ctx)

	conn := &websocket.Conn{}
	go func() {
		buildWebsocketHandler(ctx, c)(conn)
		cancel()
	}()
	for {
		select {
		case err = <-errChan:
			if !errors.Is(err, io.EOF) {
				return err
			}
		case <-done.Done():
			return nil
		}
	}
}

func TestWebsocketHandler(t *testing.T) {
	PatchConvey("Test WebsocketHandler", t, func() {
		PatchConvey("on success", func() {
			WebsocketHandlerNormalize()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldBeNil)
		})

		PatchConvey("header with no token", func() {
			Mock((*app.RequestContext).GetHeader).To(func(key string) []byte {
				return []byte("")
			}).Build()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldNotBeNil)
		})

		PatchConvey("on write ping message error", func() {
			Mock(writePingMessage).Return(errors.New("write ping message error")).Build()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldNotBeNil)
		})

		PatchConvey("on login error", func() {
			Mock(service.Core.Login).Return(errors.New("login error")).Build()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldNotBeNil)
		})

		PatchConvey("on accept error", func() {
			Mock(service.Core.Accept).Return(errors.New("accept error")).Build()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldNotBeNil)
		})

		PatchConvey("on logout error", func() {
			Mock(service.Core.Logout).Return(errors.New("logout error")).Build()
			err := CatchErrInWebsocketConn(context.Background(), &app.RequestContext{})
			So(err, ShouldNotBeNil)
		})
	})
}
