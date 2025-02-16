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

package http

import (
	"context"
	"errors"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app/client"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDo(t *testing.T) {
	cli := NewClient(&ClientConfig{})
	PatchConvey("Client.do", t, func() {
		PatchConvey("normal", func() {
			Mock((*client.Client).Do).Return(nil).Build()

			err := cli.do(context.Background(), nil, nil)
			So(err, ShouldBeNil)
		})

		PatchConvey("error", func() {
			Mock((*client.Client).Do).Return(errors.New("")).Build()

			err := cli.do(context.Background(), nil, nil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestBuildUrl(t *testing.T) {
	Convey("TestBuildUrl", t, func() {
		PatchConvey("prev[] sub[/]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost"})
			url := cli.buildUrl("/ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[/] sub[/]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/"})
			url := cli.buildUrl("/ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[] sub[]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost"})
			url := cli.buildUrl("ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[/] sub[]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/"})
			url := cli.buildUrl("ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[x] sub[./]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/"})
			url := cli.buildUrl("./ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[x] sub[../]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/"})
			url := cli.buildUrl("../ping")
			So(url, ShouldEqual, "http://localhost/ping")
		})

		PatchConvey("prev[x] sub[../../]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/a/b/c/d"})
			url := cli.buildUrl("../../ping")
			So(url, ShouldEqual, "http://localhost/a/b/ping")
		})

		PatchConvey("prev[/1/2/3] sub[../../]", func() {
			cli := NewClient(&ClientConfig{BaseUrl: "http://localhost/1/2/3"})
			url := cli.buildUrl("../..")
			So(url, ShouldEqual, "http://localhost/1")
		})
	})
}
