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
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"

	"github.com/west2-online/DomTok/app/order/domain/repository"
	"github.com/west2-online/DomTok/pkg/constants"
	"github.com/west2-online/DomTok/pkg/errno"
)

type locker struct {
	rs      *redsync.Redsync
	locks   map[int64]*redsync.Mutex
	expires map[int64]time.Time
	mu      sync.Mutex
}

func NewLocker(rs *redsync.Redsync) repository.Locker {
	l := &locker{
		rs:      rs,
		locks:   make(map[int64]*redsync.Mutex),
		expires: make(map[int64]time.Time),
	}
	l.init()
	return l
}

func (l *locker) init() {
	go func() {
		time.Sleep(constants.OrderRedSyncDefaultInterval)
		l.mu.Lock()
		now := time.Now().UnixMilli()
		for orderId, expire := range l.expires {
			if now > expire.UnixMilli() {
				delete(l.locks, orderId)
				delete(l.expires, orderId)
			}
		}
		l.mu.Unlock()
	}()
}

func (l *locker) LockOrder(orderID int64) error {
	l.mu.Lock()
	mutex, ok := l.locks[orderID]
	if !ok {
		mutex = l.rs.NewMutex(getKey(orderID))
		l.locks[orderID] = mutex
	}
	l.mu.Unlock()

	if err := mutex.Lock(); err != nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("failed to lock order: %v", err))
	}
	l.expires[orderID] = time.Now().Add(constants.OrderRedSyncDefaultTTL)
	return nil
}

func (l *locker) UnlockOrder(orderID int64) (err error) {
	l.mu.Lock()
	mutex, ok := l.locks[orderID]
	l.mu.Unlock()

	if !ok {
		return errno.NewErrNo(errno.InternalServiceErrorCode, "order not locked")
	}

	if ok, err = mutex.Unlock(); !ok || err != nil {
		return errno.NewErrNo(errno.InternalServiceErrorCode, fmt.Sprintf("failed to unlock order,err: %v", err))
	}
	return nil
}

func getKey(orderID int64) string {
	return fmt.Sprintf(constants.OrderLockFormat, orderID)
}
