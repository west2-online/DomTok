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

package utils

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch             int64 = 1577808000000 // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     int64 = 41            // 时间戳占用位数
	datacenteridBits  int64 = 5             // 数据中心id所占位数
	workeridBits      int64 = 5             // 机器id所占位数
	sequenceBits      int64 = 12            // 序列所占的位数
	timestampMax      int64 = 2199023255551 // (-1 ^ (-1 << timestampBits))时间戳最大值
	datacenteridMax   int64 = 31            // (-1 ^ (-1 << datacenteridBits))支持的最大数据中心id数量
	workeridMax       int64 = 31            // (-1 ^ (-1 << workeridBits)) 支持的最大机器id数量
	sequenceMask      int64 = 4095          // (-1 ^ (-1 << sequenceBits))支持的最大序列id数量
	workeridShift     int64 = 12            // sequenceBits 机器id左移位数
	datacenteridShift int64 = 17            // sequenceBits + workeridBits 数据中心id左移位数
	timestampShift    int64 = 22            // sequenceBits + workeridBits + datacenteridBits 时间戳左移位数

	NanosecondsInAMillisecond = 1_000_000 // 每毫秒的纳秒数
	MillisecondsInASecond     = 1000      // 每秒的毫秒数
)

type Snowflake struct {
	sync.Mutex
	timestamp    int64
	workerid     int64
	datacenterid int64
	sequence     int64
}

func NewSnowflake(datacenterid, workerid int64) (*Snowflake, error) {
	if datacenterid < 0 || datacenterid > datacenteridMax {
		return nil, fmt.Errorf("datacenterid must be between 0 and %d", datacenteridMax-1)
	}
	if workerid < 0 || workerid > workeridMax {
		return nil, fmt.Errorf("workerid must be between 0 and %d", workeridMax-1)
	}
	return &Snowflake{
		timestamp:    0,
		datacenterid: datacenterid,
		workerid:     workerid,
		sequence:     0,
	}, nil
}

// timestamp + 数据中心id + 工作节点id + 自旋id
func (s *Snowflake) NextVal() (int64, error) {
	s.Lock()
	now := time.Now().UnixNano() / NanosecondsInAMillisecond // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / NanosecondsInAMillisecond
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return 0, fmt.Errorf("epoch must be between 0 and %d", timestampMax-1)
	}
	s.timestamp = now
	r := (t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence)
	s.Unlock()
	return r, nil
}

// GetDeviceID 获取数据中心ID和机器ID
func GetDeviceID(sid int64) (datacenterid, workerid int64) {
	datacenterid = (sid >> datacenteridShift) & datacenteridMax
	workerid = (sid >> workeridShift) & workeridMax
	return
}

// GetTimestamp 获取时间戳
func GetTimestamp(sid int64) (timestamp int64) {
	timestamp = (sid >> timestampShift) & timestampMax
	return
}

// GetGenTimestamp 获取创建ID时的时间戳
func GetGenTimestamp(sid int64) (timestamp int64) {
	timestamp = GetTimestamp(sid) + epoch
	return
}

// GetGenTime 获取创建ID时的时间字符串(精度：秒)
func GetGenTime(sid int64) (t string) {
	// 需将GetGenTimestamp获取的时间戳/1000转换成秒
	t = time.Unix(GetGenTimestamp(sid)/MillisecondsInASecond, 0).Format("2006-01-02 15:04:05")
	return
}

// GetTimestampStatus 获取时间戳已使用的占比：范围（0.0 - 1.0）
func GetTimestampStatus() (state float64) {
	state = float64((time.Now().UnixNano()/NanosecondsInAMillisecond - epoch)) / float64(timestampMax)
	return
}
