// Copyright 2019 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv1

import (
	"fmt"
	"log/slog"
	"sync/atomic"
)

var logger atomic.Pointer[slog.Logger]

func init() { logger.Store(slog.Default()) }

// SetLogger sets the logger used for informational background messages.
// A nil logger disables logging. It is safe to call concurrently with logging.
func SetLogger(l *slog.Logger) { logger.Store(l) }

func logf(format string, v ...interface{}) {
	if l := logger.Load(); l != nil {
		l.Info(fmt.Sprintf(format, v...))
	}
}
