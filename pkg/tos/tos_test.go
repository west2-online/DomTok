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

package tos

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/config"
)

func isSetAccessKey() bool {
	return os.Getenv("TOS_ACCESS_KEY") != ""
}

func mockTosConfig() {
	mockey.MockValue(&config.Tos).To(&config.TosConfig{
		AccessKey: os.Getenv("TOS_ACCESS_KEY"),
		SecretKey: os.Getenv("TOS_SECRET_KEY"),
		Bucket:    "domtok",
		Region:    "cn-beijing",
		Endpoint:  "tos-cn-beijing.volces.com",
	}).Patch()
	fmt.Println(config.Tos.AccessKey)
	fmt.Println(config.Tos.SecretKey)
}

func TestTos(t *testing.T) {
	if !isSetAccessKey() {
		return
	}
	mockey.PatchConvey("", t, func() {
		mockTosConfig()
		defer mockey.UnPatchAll()

		ctx := context.Background()
		err := UploadImage(ctx, []byte("test"), "test.txt")
		convey.So(err, convey.ShouldBeNil)

		signedUrl, err := SignedURL("test.txt")
		convey.So(err, convey.ShouldBeNil)
		t.Logf("signedUrl: %s", signedUrl)

		err = DeleteImage(ctx, "test.txt")
		convey.So(err, convey.ShouldBeNil)
	})
}
