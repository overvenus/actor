// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package actor

import (
	"errors"
	"log"
	"os"
)

var logger = Logger(log.New(os.Stderr, "[actor] ", log.Ldate|log.Ltime|log.Lshortfile))

// Logger is used to log error messages.
type Logger interface {
	Printf(format string, v ...any)
	Panicf(format string, v ...any)
}

// SetLogger is used to set the logger for error message.
// The initial logger is os.Stderr.
func SetLogger(l Logger) error {
	if l == nil {
		return errors.New("logger is nil")
	}
	logger = l
	return nil
}
