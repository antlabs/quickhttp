package quickhttp

// Copyright 2021-2023 antlabs. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"sync"
)

// 生成的大小分别是
// 1 * 1024
// 2 * 1024
// 3 * 1024
// 4 * 1024
// 5 * 1024
// 6 * 1024
// 7 * 1024
func init() {
	for i := 1; i <= maxIndex; i++ {
		j := i
		pools = append(pools, sync.Pool{
			New: func() interface{} {
				buf := make([]byte, j*poolPage)
				return &buf
			},
		})
	}
}

const (
	poolPage = 1024
	maxIndex = 64
)

var pools = make([]sync.Pool, 0, maxIndex)
var upgradeRespPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 256)
		return &buf
	},
}

func selectIndex(n int) int {
	index := n / poolPage
	return index
}

func GetAndResetBytes(n int) (rv *[]byte) {
	rv = GetBytes(n)
	*rv = (*rv)[:0]
	return rv
}

func GetBytes(n int) (rv *[]byte) {
	if n <= poolPage {
		rv = pools[0].Get().(*[]byte)
		*rv = (*rv)[:cap(*rv)]
		return rv
	}

	index := selectIndex(n - 1)
	if index >= len(pools) {
		rv := make([]byte, n)
		return &rv
	}

	rv = pools[index].Get().(*[]byte)
	*rv = (*rv)[:cap(*rv)]
	return rv
}

// PutBytes可以接受 不是GetBytes分配出来的内存块
// 如果不是GetBytes分配出来内存块就移到向下一级移去，每一个索引的数据都是>= page * i + enum.MaxFrameHeaderSize，这样保证取出来的数据是够用的，不会太小。
func PutBytes(bytes *[]byte) {
	if cap(*bytes) == 0 {
		return
	}
	if cap(*bytes) < poolPage {
		return
	}

	newLen := cap(*bytes) - 1
	index := selectIndex(newLen)
	if (cap(*bytes))%poolPage != 0 {
		index--
		if index < 0 {
			return
		}
	}
	if index >= len(pools) {
		return
	}
	pools[index].Put(bytes)
}

func GetUpgradeRespBytes() *[]byte {
	return upgradeRespPool.Get().(*[]byte)
}

func PutUpgradeRespBytes(bytes *[]byte) {
	upgradeRespPool.Put(bytes)
}
