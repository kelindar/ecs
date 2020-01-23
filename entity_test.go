// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerial(t *testing.T) {
	s := next()
	assert.NotZero(t, s)
}
