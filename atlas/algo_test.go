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
	)

	distC := boxDistCalc(
		[2]int32{170 - 360, 33}, [2]int32{170 - 360, 33},
		[2]int32{-170, 33}, [2]int32{-170, 33},
	)
	if distA < distC {
		t.Fatalf("unexpected results")
	}
}

func BenchmarkBox(b *testing.B) {
	f := Box([2]int32{170, 33}, [2]int32{170, 33},
		false,
	)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f([2]int32{-170, 33}, [2]int32{-170, 33}, nil, false)
	}
}
