// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package gtpv1 provides GTP-U connections, userland relaying and optional
// Linux kernel GTP-U integration. Its message and ie subpackages implement
// GTPv1-C and GTPv1-U codecs following TS 29.060 V19.0.0, TS 29.281 V19.2.0
// and TS 38.415 V19.1.0. GTPv1-C has no connection layer.
//
// DialUPlane performs an Echo handshake and starts serving; Close stops the
// returned connection. NewUPlaneConn with ListenAndServe starts a server.
// Handlers run on the serving goroutine in receive order and should not block.
// Parsed values may reference the input buffer; copy data retained by handlers.
// IE constructors carry errors available through IE.Err, and accessors on nil
// IEs return ie.ErrIENotFound. SetLogger accepts *slog.Logger; nil disables
// background logging. Applications own routing and TUN/TAP integration.
package gtpv1
