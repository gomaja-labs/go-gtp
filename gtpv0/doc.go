// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package gtpv0 documents the GTPv0 codecs in its message and ie subpackages,
// following GSM 09.60 V7.10.0. It provides no connection layer or PDP context
// state machine. Use message.Marshal and message.Parse with application-owned
// datagram I/O. IE constructors carry errors available through IE.Err;
// accessors on nil IEs return ie.ErrIENotFound. Parsed values may reference
// the input buffer and must not outlive its reuse or modification.
package gtpv0
