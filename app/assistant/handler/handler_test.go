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
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/assistant/service"
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
