// Copyright (c) 2018 Josh Baker
// Copyright (c) 2020 Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
// This is a fork of https://github.com/tidwall/rbang adapted for int32 coordinate system.

package atlas

// DistanceFunc represents a distance function
type DistanceFunc = func(min, max [2]int32, data interface{}, item bool) (dist int32)

// Box performs simple box-distance algorithm on rectangles. When wrapX
// is provided, the operation does a cylinder wrapping of the X value to allow
// for antimeridian calculations. When itemDist is provided (not nil), it
// becomes the caller's responsibility to return the box-distance.
func Box(targetMin, targetMax [2]int32, wrapX bool) DistanceFunc {
	return func(min, max [2]int32, data interface{}, item bool) (dist int32) {
		return boxDistCalc(targetMin, targetMax, min, max)
	}
}

func mmin(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func mmax(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

// boxDistCalc returns the distance from rectangle A to rectangle B. When wrapX
// is provided, the operation does a cylinder wrapping of the X value to allow
// for antimeridian calculations.
func boxDistCalc(aMin, aMax, bMin, bMax [2]int32) int32 {
	var dist, squared int32

	// X
	squared = mmax(aMin[0], bMin[0]) - mmin(aMax[0], bMax[0])
	if squared > 0 {
		dist += squared * squared
	}

	// Y
	squared = mmax(aMin[1], bMin[1]) - mmin(aMax[1], bMax[1])
	if squared > 0 {
		dist += squared * squared
	}

	return dist
}
