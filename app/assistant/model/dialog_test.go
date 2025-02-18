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

package model

import (
	"sync"
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDialog(t *testing.T) {
	PatchConvey("Test Dialog", t, func() {
		PatchConvey("Test IDialog Impl", func() {
			d := NewDialog("1", "Hello")
			So(d.Unique(), ShouldEqual, "1")
			So(d.Message(), ShouldEqual, "Hello")
		})

		PatchConvey("Test Send", func() {
			d := NewDialog("1", "Hello")
			go func() {
				d.Send("World")
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case msg := <-d._receiver:
				So(msg, ShouldEqual, "World")
			case <-ticker.C:
				t.Error("timeout")
			}
		})

		PatchConvey("Test Close", func() {
			d := NewDialog("1", "Hello")
			go func() {
				d.Close()
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case <-d.ctx.Done():
			case <-ticker.C:
				t.Error("timeout")
			}
		})

		PatchConvey("Test NotifyOnMessage", func() {
			d := NewDialog("1", "Hello")
			go func() {
				d.Send("World")
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case msg := <-d.NotifyOnMessage():
				So(msg, ShouldEqual, "World")
			case <-ticker.C:
				t.Error("timeout")
			}
		})

		PatchConvey("Test NotifyOnClosed", func() {
			d := NewDialog("1", "Hello")
			go func() {
				d.Close()
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case <-d.NotifyOnClosed():
			case <-ticker.C:
				t.Error("timeout")
			}
		})

		PatchConvey("Test Close After Send", func() {
			d := NewDialog("1", "Hello")
			isStillBlocked := true
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				d.Send("World")
				d.Close()
				isStillBlocked = false
				wg.Done()
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case msg := <-d.NotifyOnMessage():
				So(msg, ShouldEqual, "World")
			case <-d.NotifyOnClosed():
			case <-ticker.C:
				wg.Done()
				return
			}
			wg.Wait()
			So(isStillBlocked, ShouldBeFalse)
		})

		PatchConvey("Test Send Many Times After Close", func() {
			d := NewDialog("1", "Hello")
			d.Close()
			go func() {
				d.Send("World")
				d.Send("World")
				d.Send("World")
			}()
			ticker := time.NewTicker(1 * time.Second)
			select {
			case <-d.NotifyOnClosed():
			case <-ticker.C:
				t.Error("timeout")
			}
		})
	})
}
