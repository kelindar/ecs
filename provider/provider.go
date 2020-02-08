// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package provider

//go:generate genny -pkg=provider -in=generic.go -out=z_components.go gen "Any=Point"
//go:generate genny -pkg=provider -in=generic_test.go -out=z_components_test.go gen "Any=Point"

// Point represents a 3D point component, provided as an example
type Point struct {
	X, Y, Z int32
}
