// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package gtpv2 provides GTPv2-C connections, handlers, transactions and
// in-memory sessions and bearers following TS 29.274 V20.0.0. Its message
// and ie subpackages provide typed codecs.
//
// NewConn with ListenAndServe starts a server. Dial performs an Echo handshake
// and starts serving; its context bounds only the handshake. Close stops the
// returned connection. Register handlers before serving. Handlers run on the
// serving goroutine in receive order; do not wait in a handler for a response
// that the same receive loop must dispatch. Applications implement procedure
// policy, peer authorization and persistence.
//
// Decoding checks framing and applies message-specific IE content error rules;
// callers must inspect errors before using decoded values. IE constructors
// carry errors available through IE.Err, and accessors on nil IEs return
// ie.ErrIENotFound. Parsed values may reference the input buffer; copy data
// retained by handlers. SetLogger accepts *slog.Logger; nil disables logging.
package gtpv2
