//go:build !goexperiment.simd

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

// When SIMD experiment is not enabled, all SIMD functions return nil
// to indicate they are not available. The caller will use fallback implementations.

func vectorAddSIMD(c Context, u, v *Vector) *Vector {
	return nil
}

func vectorSubSIMD(c Context, u, v *Vector) *Vector {
	return nil
}

func vectorMulSIMD(c Context, u, v *Vector) *Vector {
	return nil
}

func vectorMinSIMD(c Context, u, v *Vector) *Vector {
	return nil
}

func vectorMaxSIMD(c Context, u, v *Vector) *Vector {
	return nil
}
