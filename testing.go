// Copyright 2021 PingCAP, Inc.
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

// NewRouter4Test returns a new router. Test only.
func NewRouter4Test[T any](name string) *Router[T] {
	return newRouter[T](name, defaultRouterChunkCap)
}

// InsertMailbox4Test add a mailbox into router. Test only.
func (r *Router[T]) InsertMailbox4Test(id ID, mb Mailbox[T]) {
	r.procs.Store(id, &proc[T]{mb: mb})
}
