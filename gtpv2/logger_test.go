// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package gtpv2

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestLoggerReconfiguration(t *testing.T) {
	defer SetLogger(slog.Default())
	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, nil))
	SetLogger(l)
	logf("value %d", 42)
	if !strings.Contains(buf.String(), "value 42") {
		t.Fatal(buf.String())
	}
	before := buf.Len()
	SetLogger(nil)
	logf("disabled")
	if buf.Len() != before {
		t.Fatal("nil logger did not disable logging")
	}
	SetLogger(l)
	logf("enabled")
	if !strings.Contains(buf.String(), "enabled") {
		t.Fatal("logging did not resume")
	}
}

func TestLoggerConcurrent(t *testing.T) {
	defer SetLogger(slog.Default())
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	var wg sync.WaitGroup
	for n := 0; n < 8; n++ {
		wg.Go(func() {
			for k := 0; k < 1000; k++ {
				SetLogger(l)
				logf("value %d", k)
				SetLogger(nil)
				logf("disabled")
			}
		})
	}
	wg.Wait()
}
