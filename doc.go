// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package gtp marshals and parses GTPv0, GTPv1 and GTPv2 messages.
// Parse dispatches by the wire version; version-specific message and ie
// subpackages provide typed codecs. Parsed values may reference the input
// buffer, so copy it before parsing when values must survive buffer reuse.
//
// Package gtpv1 provides userland and optional Linux kernel GTP-U connections;
// gtpv2 provides GTPv2-C connections, handlers, transactions and in-memory
// sessions. GTPv0 and GTPv1-C have codecs only. See the version packages for
// decode error policy and the README for supported specification editions.
package gtp
