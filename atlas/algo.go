// Copyright (c) 2018 Josh Baker
// Copyright (c) 2020 Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.
// This is a fork of https://github.com/tidwall/rbang adapted for int32 coordinate system.

package atlas

// Box performs simple box-distance algorithm on rectangles. When wrapX
// is provided, the operation does a cylinder wrapping of the X value to allow
// for antimeridian calculations. When itemDist is provided (not nil), it
// becomes the caller's responsibility to return the box-distance.
func Box(
	targetMin, targetMax [2]int32, wrapX bool,
	itemDist func(min, max [2]int32, data interface{}) int32,
) (
	algo func(min, max [2]int32, data interface{}, item bool) (dist int32),
) {
	return func(min, max [2]int32, data interface{}, item bool) (dist int32) {
		if item && itemDist != nil {
			return itemDist(min, max, data)
		}
		return boxDistCalc(targetMin, targetMax, min, max, wrapX)
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
func boxDistCalc(aMin, aMax, bMin, bMax [2]int32, wrapX bool) int32 {
	var dist int32
	var squared int32

	// X
	squared = mmax(aMin[0], bMin[0]) - mmin(aMax[0], bMax[0])
	if wrapX {
		squaredLeft := mmax(aMin[0]-360, bMin[0]) - mmin(aMax[0]-360, bMax[0])
		squaredRight := mmax(aMin[0]+360, bMin[0]) - mmin(aMax[0]+360, bMax[0])
		squared = mmin(squared, mmin(squaredLeft, squaredRight))
	}
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

// contains return struct when b is fully contained inside of n
func intersects(aMin, aMax, bMin, bMax [2]int32) bool {
	if bMin[0] > aMax[0] || bMax[0] < aMin[0] {
		return false
	}
	if bMin[1] > aMax[1] || bMax[1] < aMin[1] {
		return false
	}
	return true
}
