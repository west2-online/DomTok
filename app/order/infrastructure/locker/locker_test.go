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

package locker

import (
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/config"
	"github.com/west2-online/DomTok/pkg/base/client"
	"github.com/west2-online/DomTok/pkg/logger"
	"github.com/west2-online/DomTok/pkg/utils"
)

func initLocker(t *testing.T) repository.Locker {
	t.Helper()
	config.Init("test")
	logger.Ignore()
	c, err := client.InitRedis(0)
	if err != nil {
		t.Fatalf("failed to init redis client: %v", err)
	}
	rs := client.InitRedSync(c)
	return NewLocker(rs)
}

func TestOrder_Locker(t *testing.T) {
	if !utils.EnvironmentEnable() {
		return
	}
	l := initLocker(t)

	Convey("Test two threads acquire locks at the same time", t, func() {
		Convey("Normal get lock", func() {
			id := rand.Int64()
			So(l.LockOrder(id), ShouldBeNil)
			So(l.UnlockOrder(id), ShouldBeNil)
		})

		Convey("Test repeat unlock", func() {
			id := rand.Int64()
			So(l.LockOrder(id), ShouldBeNil)
			So(l.UnlockOrder(id), ShouldBeNil)
			So(l.UnlockOrder(id), ShouldNotBeNil)
		})

		Convey("Test unlock with no lock", func() {
			id := rand.Int64()
			So(l.UnlockOrder(id), ShouldNotBeNil)
		})

		Convey("Test several goroutine", func() {
			var wg sync.WaitGroup
			Convey("Test 2 goroutine acquire the same id", func() {
				acquired := make(chan struct{}) // 用于通知锁已被获取
				proceed := make(chan struct{})  // 用于通知释放锁
				id := rand.Int64()
				wg.Add(2)
				var getLock bool

				// 第一个 goroutine：获取锁并阻塞直到收到释放信号
				go func() {
					defer wg.Done()
					_ = l.LockOrder(id)
					close(acquired) // 通知锁已获取
					<-proceed       // 等待释放信号
					_ = l.UnlockOrder(id)
				}()

				// 第二个 goroutine：等待锁被获取后尝试获取锁
				go func() {
					defer wg.Done()
					<-acquired // 确保第一个已持有锁

					// 尝试非阻塞获取锁，预期失败
					lc := l.(*locker) //nolint
					if lc.rs.NewMutex(getKey(id)).TryLock() == nil {
						// 如果真的拿到锁了 (预期中不应该走到这)
						getLock = true
					}
					close(proceed) // 允许第一个 goroutine 释放锁
				}()

				// 等待两个 goroutine 就绪
				wg.Wait()
				So(getLock, ShouldBeFalse)
			})

			Convey("Test several threads use different goroutines", func() {
				id := rand.Int64()
				times := 5
				wg.Add(times)
				getLock := atomic.Int32{}
				fn := func(id int64, index int64) {
					defer wg.Done()
					_ = l.LockOrder(id + index)
					getLock.Add(1)
					_ = l.UnlockOrder(id + index)
				}
				for i := 0; i < times; i++ {
					go fn(id, int64(i))
				}
				wg.Wait()
				So(getLock.Load(), ShouldEqual, times)
			})
		})
	})
}
