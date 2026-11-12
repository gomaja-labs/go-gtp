// Copyright 2019-2024 go-gtp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-gtp/gtpv0/ie"
	"github.com/gomaja-labs/go-gtp/gtpv0/message"
	"github.com/gomaja-labs/go-gtp/utils"
)

func TestMessagePropagatesConstructionError(t *testing.T) {
	i := ie.NewMSISDN("12x")
	m := message.NewCreatePDPContextRequest(0, 0, 0, i)
	if _, err := m.Marshal(); !errors.Is(err, utils.ErrMalformed) {
		t.Fatalf("construction error lost: %v", err)
	}
}
