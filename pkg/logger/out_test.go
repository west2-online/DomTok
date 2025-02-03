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

package logger

import (
	"os"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_checkAndOpenFile(t *testing.T) {
	Convey("Test checkAndOpenFile", t, func() {
		path := "./test.log"
		f := checkAndOpenFile(path)
		So(f.Name(), ShouldEqual, path)
		fi, err := os.Stat(path)
		So(err, ShouldBeNil)
		So(fi.Name(), ShouldEqual, "test.log")

		// 触发 Finalizer
		f = nil
		runtime.GC()

		So(os.Remove(path), ShouldBeNil)
	})
}
