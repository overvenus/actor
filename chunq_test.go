// Copyright 2023 Neil Shen.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package actor

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type qModel struct {
	queue *list.List
}

func (q *qModel) Len() int {
	return q.queue.Len()
}

func (q *qModel) PushBack(v int) {
	q.queue.PushBack(v)
}

func (q *qModel) PopFront() (int, bool) {
	if q.queue.Front() == nil {
		return 0, false
	}
	front := q.queue.Front()
	q.queue.Remove(front)
	return front.Value.(int), true
}

func TestChunkQ(t *testing.T) {
	t.Parallel()

	seed := time.Now().Unix()
	// seed := int64(1673600359)
	rnd := rand.New(rand.NewSource(seed))
	model := qModel{queue: list.New()}

	q := newChunkQ[int](17)
	require.EqualValues(t, 0, q.Len())

	for init := 0; init < 128; init++ {
		for i := 0; i < init; i++ {
			q.PushBack(i)
			model.PushBack(i)
		}
		for i := init; i < init+200; i++ {
			op := rnd.Uint32() & (math.MaxUint32 >> 31)
			if op == 0 {
				q.PushBack(i)
				model.PushBack(i)
				t.Logf("push %d\n", i)
			} else if op == 1 {
				v1, ok1 := q.PopFront()
				v2, ok2 := model.PopFront()
				t.Logf("pop (%d, %v) (%d, %v)\n", v1, ok1, v2, ok2)
				require.Equal(t, v2, v1, "%d %d", seed, init)
				require.Equal(t, ok2, ok1, "%d %d", seed, init)
			} else {
				t.Fatal("op should be 1 or 0, ", op, "%d %d", seed, init)
			}
			require.Equal(t, model.Len(), q.Len(), "%d %d", seed, init)
		}
	}
}

func TestChunkQRetainUndrainedChunk(t *testing.T) {
	t.Parallel()

	q := newChunkQ[int](17)
	q.PushBack(1)
	chk := q.queue.Front().Value.(*chunk[int])

	v, ok := q.PopFront()
	require.Equal(t, 1, v)
	require.True(t, ok)

	// The chunk is retain.
	require.Equal(t, 1, q.queue.Len())
	require.Equal(t, chk, q.queue.Front().Value.(*chunk[int]))

	// PopFront again.
	_, ok = q.PopFront()
	require.False(t, ok)
	// The chunk is retain.
	require.Equal(t, 1, q.queue.Len())
	require.Equal(t, chk, q.queue.Front().Value.(*chunk[int]))
}

func BenchmarkQueue(b *testing.B) {
	pushPopListQ := qModel{queue: list.New()}
	b.Run("list_queue", func(b *testing.B) {
		benchmarkQueue(b, func(i int) {
			q := qModel{queue: list.New()}
			for j := 0; j < i; j++ {
				q.PushBack(j)
			}
		}, func(i int) {
			for j := 0; j < i; j++ {
				pushPopListQ.PushBack(j)
			}
			for j := 0; j < i; j++ {
				pushPopListQ.PopFront()
			}
		})
	})

	pushPopChunkQ := newChunkQ[int](64)
	b.Run("chunk_queue_cap_64", func(b *testing.B) {
		benchmarkQueue(b, func(i int) {
			q := newChunkQ[int](64)
			for j := 0; j < i; j++ {
				q.PushBack(j)
			}
		}, func(i int) {
			for j := 0; j < i; j++ {
				pushPopChunkQ.PushBack(j)
			}
			for j := 0; j < i; j++ {
				pushPopChunkQ.PopFront()
			}
		})
	})
}

func benchmarkQueue(b *testing.B, push, pushPop func(int)) {
	for item := 32; item <= 1024; item = item * 2 {
		b.Run(fmt.Sprintf("push_%d_item", item), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				push(item)
			}
		})

		b.Run(fmt.Sprintf("push_pop_%d_item", item), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pushPop(item)
			}
		})
	}
}
