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
	"sync"
)

type chunk[T any] struct {
	buffer []T
	front  int
	back   int
}

func newChunk[T any](cap int) *chunk[T] {
	return &chunk[T]{
		buffer: make([]T, cap),
	}
}

func (c *chunk[T]) reset() {
	c.front, c.back = 0, 0
}

func (c *chunk[T]) isDrained() bool {
	return cap(c.buffer) == c.back
}

func (c *chunk[T]) push(v T) bool {
	if c.back == cap(c.buffer) {
		return false
	}
	c.buffer[c.back] = v
	c.back++
	return true
}

func (c *chunk[T]) pop() (T, bool) {
	var empty T
	if c.front == c.back {
		return empty, false
	}
	c.front++
	v := c.buffer[c.front-1]
	c.buffer[c.front-1] = empty // Do not hold reference, prevent memory leak.
	return v, true
}

type chunkQ[T any] struct {
	cache    *sync.Pool
	queue    *list.List
	chunkCap int
	len      int
}

func newChunkQ[T any](chunkCap int) *chunkQ[T] {
	return &chunkQ[T]{
		cache: &sync.Pool{New: func() any {
			return newChunk[T](chunkCap)
		}},
		queue:    list.New(),
		chunkCap: chunkCap,
	}
}

func (q *chunkQ[T]) Len() int {
	return q.len
}

func (q *chunkQ[T]) PushBack(v T) {
	if q.queue.Back() == nil {
		q.queue.PushBack(q.cache.Get())
	}
	ok := q.queue.Back().Value.(*chunk[T]).push(v)
	if !ok {
		chunk := q.cache.Get().(*chunk[T])
		chunk.push(v)
		q.queue.PushBack(chunk)
	}
	q.len++
}

func (q *chunkQ[T]) PopFront() (T, bool) {
	for front := q.queue.Front(); front != nil; front = q.queue.Front() {
		chunk := front.Value.(*chunk[T])
		v, ok := chunk.pop()
		if !ok {
			if chunk.isDrained() {
				// Only remove chunk when all its buffer was used.
				q.queue.Remove(front)
				// Recycle chunk to save memory allocation.
				chunk.reset()
				q.cache.Put(chunk)
				continue
			}
			// The chunk is no drained, it can be used in later push.
			break
		}
		q.len--
		return v, ok
	}
	var v T
	return v, false
}
