// Copyright (c) 2018 Josh Baker
// Copyright (c) 2020 Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
// This is a fork of https://github.com/tidwall/rbang adapted for int32 coordinate system.

package atlas

import (
	"testing"
)

func TestBoxDist(t *testing.T) {
	distA := boxDistCalc(
		[2]int32{170, 33}, [2]int32{170, 33},
		[2]int32{-170, 33}, [2]int32{-170, 33},
		false,
	)

	distB := boxDistCalc(
		[2]int32{170, 33}, [2]int32{170, 33},
		[2]int32{-170, 33}, [2]int32{-170, 33},
		true,
	)
	distC := boxDistCalc(
		[2]int32{170 - 360, 33}, [2]int32{170 - 360, 33},
		[2]int32{-170, 33}, [2]int32{-170, 33},
		false,
	)
	if distA < distB || distC != distB {
		t.Fatalf("unexpected results")
	}
}
