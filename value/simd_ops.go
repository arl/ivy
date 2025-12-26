//go:build goexperiment.simd

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"simd/archsimd"
)

// This file implements vector operations using Go's experimental SIMD API.
// It requires building with GOEXPERIMENT=simd.
//
// The SIMD API provides explicit SIMD vector types (Int64x2, Int64x4, Int64x8)
// that map directly to hardware SIMD registers (SSE, AVX2, AVX512).

// vectorAddSIMD performs element-wise addition using SIMD instructions.
// It processes chunks of the vector using Int64x4 (256-bit AVX2) vectors.
func vectorAddSIMD(c Context, u, v *Vector) *Vector {
	if u.Len() != v.Len() {
		Errorf("length mismatch in vector addition")
	}
	n := u.Len()
	if n == 0 {
		return u
	}
	
	// Check if both vectors contain only small Ints
	if !u.AllInts() || !v.AllInts() {
		return nil // Fallback to generic implementation
	}
	
	result := newVectorEditor(n, nil)
	
	// Check if AVX2 is available (for Int64x4)
	if archsimd.X86.HasAVX2() && n >= 4 {
		// Process 4 elements at a time using Int64x4 (AVX2)
		i := 0
		for i+4 <= n {
			// Load 4 int64 values into SIMD vectors
			var uVals, vVals [4]int64
			for j := 0; j < 4; j++ {
				uVals[j] = int64(u.At(i+j).(Int))
				vVals[j] = int64(v.At(i+j).(Int))
			}
			
			// Load into SIMD vectors
			uVec := archsimd.LoadInt64x4(&uVals)
			vVec := archsimd.LoadInt64x4(&vVals)
			
			// Perform SIMD addition
			resVec := uVec.Add(vVec)
			
			// Store results back
			var resVals [4]int64
			resVec.Store(&resVals)
			
			for j := 0; j < 4; j++ {
				result.Set(i+j, Int(resVals[j]).maybeBig())
			}
			i += 4
		}
		
		// Handle remaining elements
		for ; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui + vi).maybeBig())
		}
	} else {
		// Fallback to scalar processing
		for i := 0; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui + vi).maybeBig())
		}
	}
	
	return result.Publish()
}

// vectorSubSIMD performs element-wise subtraction using SIMD instructions.
func vectorSubSIMD(c Context, u, v *Vector) *Vector {
	if u.Len() != v.Len() {
		Errorf("length mismatch in vector subtraction")
	}
	n := u.Len()
	if n == 0 {
		return u
	}
	
	if !u.AllInts() || !v.AllInts() {
		return nil
	}
	
	result := newVectorEditor(n, nil)
	
	if archsimd.X86.HasAVX2() && n >= 4 {
		i := 0
		for i+4 <= n {
			var uVals, vVals [4]int64
			for j := 0; j < 4; j++ {
				uVals[j] = int64(u.At(i+j).(Int))
				vVals[j] = int64(v.At(i+j).(Int))
			}
			
			uVec := archsimd.LoadInt64x4(&uVals)
			vVec := archsimd.LoadInt64x4(&vVals)
			resVec := uVec.Sub(vVec)
			
			var resVals [4]int64
			resVec.Store(&resVals)
			
			for j := 0; j < 4; j++ {
				result.Set(i+j, Int(resVals[j]).maybeBig())
			}
			i += 4
		}
		
		for ; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui - vi).maybeBig())
		}
	} else {
		for i := 0; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui - vi).maybeBig())
		}
	}
	
	return result.Publish()
}

// vectorMulSIMD performs element-wise multiplication using SIMD instructions.
func vectorMulSIMD(c Context, u, v *Vector) *Vector {
	if u.Len() != v.Len() {
		Errorf("length mismatch in vector multiplication")
	}
	n := u.Len()
	if n == 0 {
		return u
	}
	
	if !u.AllInts() || !v.AllInts() {
		return nil
	}
	
	result := newVectorEditor(n, nil)
	
	// Note: Int64 multiplication requires AVX512DQ for VPMULLQ instruction
	if archsimd.X86.HasAVX512DQ() && n >= 8 {
		i := 0
		for i+8 <= n {
			var uVals, vVals [8]int64
			for j := 0; j < 8; j++ {
				uVals[j] = int64(u.At(i+j).(Int))
				vVals[j] = int64(v.At(i+j).(Int))
			}
			
			uVec := archsimd.LoadInt64x8(&uVals)
			vVec := archsimd.LoadInt64x8(&vVals)
			resVec := uVec.Mul(vVec)
			
			var resVals [8]int64
			resVec.Store(&resVals)
			
			for j := 0; j < 8; j++ {
				result.Set(i+j, Int(resVals[j]).maybeBig())
			}
			i += 8
		}
		
		for ; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui * vi).maybeBig())
		}
	} else {
		// Fallback to scalar (Int64 SIMD mul needs AVX512DQ)
		for i := 0; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			result.Set(i, (ui * vi).maybeBig())
		}
	}
	
	return result.Publish()
}

// vectorMinSIMD performs element-wise minimum using SIMD instructions.
func vectorMinSIMD(c Context, u, v *Vector) *Vector {
	if u.Len() != v.Len() {
		Errorf("length mismatch in vector min")
	}
	n := u.Len()
	if n == 0 {
		return u
	}
	
	if !u.AllInts() || !v.AllInts() {
		return nil
	}
	
	result := newVectorEditor(n, nil)
	
	// Int64 min requires AVX512VL or AVX512F
	if archsimd.X86.HasAVX512F() && n >= 8 {
		i := 0
		for i+8 <= n {
			var uVals, vVals [8]int64
			for j := 0; j < 8; j++ {
				uVals[j] = int64(u.At(i+j).(Int))
				vVals[j] = int64(v.At(i+j).(Int))
			}
			
			uVec := archsimd.LoadInt64x8(&uVals)
			vVec := archsimd.LoadInt64x8(&vVals)
			resVec := uVec.Min(vVec)
			
			var resVals [8]int64
			resVec.Store(&resVals)
			
			for j := 0; j < 8; j++ {
				result.Set(i+j, Int(resVals[j]))
			}
			i += 8
		}
		
		for ; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			if ui < vi {
				result.Set(i, ui)
			} else {
				result.Set(i, vi)
			}
		}
	} else {
		for i := 0; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			if ui < vi {
				result.Set(i, ui)
			} else {
				result.Set(i, vi)
			}
		}
	}
	
	return result.Publish()
}

// vectorMaxSIMD performs element-wise maximum using SIMD instructions.
func vectorMaxSIMD(c Context, u, v *Vector) *Vector {
	if u.Len() != v.Len() {
		Errorf("length mismatch in vector max")
	}
	n := u.Len()
	if n == 0 {
		return u
	}
	
	if !u.AllInts() || !v.AllInts() {
		return nil
	}
	
	result := newVectorEditor(n, nil)
	
	if archsimd.X86.HasAVX512F() && n >= 8 {
		i := 0
		for i+8 <= n {
			var uVals, vVals [8]int64
			for j := 0; j < 8; j++ {
				uVals[j] = int64(u.At(i+j).(Int))
				vVals[j] = int64(v.At(i+j).(Int))
			}
			
			uVec := archsimd.LoadInt64x8(&uVals)
			vVec := archsimd.LoadInt64x8(&vVals)
			resVec := uVec.Max(vVec)
			
			var resVals [8]int64
			resVec.Store(&resVals)
			
			for j := 0; j < 8; j++ {
				result.Set(i+j, Int(resVals[j]))
			}
			i += 8
		}
		
		for ; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			if ui > vi {
				result.Set(i, ui)
			} else {
				result.Set(i, vi)
			}
		}
	} else {
		for i := 0; i < n; i++ {
			ui := u.At(i).(Int)
			vi := v.At(i).(Int)
			if ui > vi {
				result.Set(i, ui)
			} else {
				result.Set(i, vi)
			}
		}
	}
	
	return result.Publish()
}
