// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"testing"
	
	"robpike.io/ivy/config"
)

// Test helpers

func getTestContext() Context {
	conf := &config.Config{}
	return benchContext{config: conf}
}

func vectorsEqual(t *testing.T, v1, v2 *Vector) bool {
	t.Helper()
	if v1.Len() != v2.Len() {
		t.Errorf("vector length mismatch: %d vs %d", v1.Len(), v2.Len())
		return false
	}
	for i := 0; i < v1.Len(); i++ {
		val1 := v1.At(i)
		val2 := v2.At(i)
		if val1 != val2 {
			// For more detailed comparison
			int1, ok1 := val1.(Int)
			int2, ok2 := val2.(Int)
			if !ok1 || !ok2 || int1 != int2 {
				t.Errorf("mismatch at index %d: %v vs %v", i, val1, val2)
				return false
			}
		}
	}
	return true
}

// Tests for vectorAddInt

func TestVectorAddInt_Basic(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(1, 2, 3, 4, 5)
	v := NewIntVector(10, 20, 30, 40, 50)
	
	result := vectorAddInt(c, u, v)
	if result == nil {
		t.Fatal("vectorAddInt returned nil for valid input")
	}
	
	expected := NewIntVector(11, 22, 33, 44, 55)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorAddInt result incorrect")
	}
}

func TestVectorAddInt_EmptyVectors(t *testing.T) {
	c := getTestContext()
	u := NewVector()
	v := NewVector()
	
	result := vectorAddInt(c, u, v)
	if result == nil {
		t.Fatal("vectorAddInt returned nil for empty vectors")
	}
	
	if result.Len() != 0 {
		t.Errorf("expected empty vector, got length %d", result.Len())
	}
}

func TestVectorAddInt_LargeValues(t *testing.T) {
	c := getTestContext()
	// Use values that will stay as Int (not overflow to BigInt)
	u := NewIntVector(1000000, 2000000, 3000000)
	v := NewIntVector(4000000, 5000000, 6000000)
	
	result := vectorAddInt(c, u, v)
	if result == nil {
		t.Fatal("vectorAddInt returned nil")
	}
	
	expected := NewIntVector(5000000, 7000000, 9000000)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorAddInt result incorrect for large values")
	}
}

func TestVectorAddInt_MatchesOriginal(t *testing.T) {
	c := getTestContext()
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		u := makeIntVector(size)
		v := makeIntVector(size)
		
		original := c.EvalBinary(u, "+", v)
		simd := vectorAddInt(c, u, v)
		
		if simd == nil {
			t.Fatalf("SIMD add returned nil for size %d", size)
		}
		
		originalVec := original.(*Vector)
		if !vectorsEqual(t, originalVec, simd) {
			t.Errorf("SIMD result differs from original for size %d", size)
		}
	}
}

// Tests for vectorSubInt

func TestVectorSubInt_Basic(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(10, 20, 30, 40, 50)
	v := NewIntVector(1, 2, 3, 4, 5)
	
	result := vectorSubInt(c, u, v)
	if result == nil {
		t.Fatal("vectorSubInt returned nil")
	}
	
	expected := NewIntVector(9, 18, 27, 36, 45)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorSubInt result incorrect")
	}
}

func TestVectorSubInt_NegativeResults(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(1, 2, 3)
	v := NewIntVector(10, 20, 30)
	
	result := vectorSubInt(c, u, v)
	if result == nil {
		t.Fatal("vectorSubInt returned nil")
	}
	
	expected := NewIntVector(-9, -18, -27)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorSubInt result incorrect for negative results")
	}
}

// Tests for vectorMulInt

func TestVectorMulInt_Basic(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(2, 3, 4, 5)
	v := NewIntVector(10, 10, 10, 10)
	
	result := vectorMulInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMulInt returned nil")
	}
	
	expected := NewIntVector(20, 30, 40, 50)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMulInt result incorrect")
	}
}

func TestVectorMulInt_WithZeros(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(1, 2, 3, 4)
	v := NewIntVector(0, 1, 0, 1)
	
	result := vectorMulInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMulInt returned nil")
	}
	
	expected := NewIntVector(0, 2, 0, 4)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMulInt result incorrect with zeros")
	}
}

// Tests for vectorMinInt

func TestVectorMinInt_Basic(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(1, 5, 3, 8, 2)
	v := NewIntVector(4, 2, 6, 1, 9)
	
	result := vectorMinInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMinInt returned nil")
	}
	
	expected := NewIntVector(1, 2, 3, 1, 2)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMinInt result incorrect")
	}
}

func TestVectorMinInt_NegativeValues(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(-5, -2, 0, 3)
	v := NewIntVector(-3, -6, 1, 2)
	
	result := vectorMinInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMinInt returned nil")
	}
	
	expected := NewIntVector(-5, -6, 0, 2)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMinInt result incorrect with negative values")
	}
}

// Tests for vectorMaxInt

func TestVectorMaxInt_Basic(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(1, 5, 3, 8, 2)
	v := NewIntVector(4, 2, 6, 1, 9)
	
	result := vectorMaxInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMaxInt returned nil")
	}
	
	expected := NewIntVector(4, 5, 6, 8, 9)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMaxInt result incorrect")
	}
}

func TestVectorMaxInt_NegativeValues(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(-5, -2, 0, 3)
	v := NewIntVector(-3, -6, 1, 2)
	
	result := vectorMaxInt(c, u, v)
	if result == nil {
		t.Fatal("vectorMaxInt returned nil")
	}
	
	expected := NewIntVector(-3, -2, 1, 3)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorMaxInt result incorrect with negative values")
	}
}

// Tests for scalar operations

func TestVectorScalarAddInt_Basic(t *testing.T) {
	c := getTestContext()
	scalar := Int(10)
	v := NewIntVector(1, 2, 3, 4, 5)
	
	result := vectorScalarAddInt(c, scalar, v)
	if result == nil {
		t.Fatal("vectorScalarAddInt returned nil")
	}
	
	expected := NewIntVector(11, 12, 13, 14, 15)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorScalarAddInt result incorrect")
	}
}

func TestVectorScalarMulInt_Basic(t *testing.T) {
	c := getTestContext()
	scalar := Int(3)
	v := NewIntVector(1, 2, 3, 4, 5)
	
	result := vectorScalarMulInt(c, scalar, v)
	if result == nil {
		t.Fatal("vectorScalarMulInt returned nil")
	}
	
	expected := NewIntVector(3, 6, 9, 12, 15)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorScalarMulInt result incorrect")
	}
}

// Tests for fallback to nil (non-Int vectors)

func TestVectorAddInt_ReturnsNilForNonInt(t *testing.T) {
	c := getTestContext()
	// Create a vector with BigInt
	u := NewVector(BigInt{bigInt64(1000000000000).Int})
	v := NewVector(BigInt{bigInt64(2000000000000).Int})
	
	result := vectorAddInt(c, u, v)
	if result != nil {
		t.Errorf("vectorAddInt should return nil for BigInt vectors, got %v", result)
	}
}

// Edge case tests

func TestVectorAddInt_SingleElement(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(42)
	v := NewIntVector(58)
	
	result := vectorAddInt(c, u, v)
	if result == nil {
		t.Fatal("vectorAddInt returned nil for single element")
	}
	
	expected := NewIntVector(100)
	if !vectorsEqual(t, result, expected) {
		t.Errorf("vectorAddInt result incorrect for single element")
	}
}

func TestVectorMinMax_SameValues(t *testing.T) {
	c := getTestContext()
	u := NewIntVector(5, 5, 5)
	v := NewIntVector(5, 5, 5)
	
	minResult := vectorMinInt(c, u, v)
	maxResult := vectorMaxInt(c, u, v)
	
	if minResult == nil || maxResult == nil {
		t.Fatal("min/max returned nil")
	}
	
	expected := NewIntVector(5, 5, 5)
	if !vectorsEqual(t, minResult, expected) {
		t.Errorf("min result incorrect for same values")
	}
	if !vectorsEqual(t, maxResult, expected) {
		t.Errorf("max result incorrect for same values")
	}
}

// Integration tests

func TestSIMDOperations_Integration(t *testing.T) {
	c := getTestContext()
	
	// Test that SIMD operations integrate correctly with the binary op system
	u := NewIntVector(1, 2, 3)
	v := NewIntVector(4, 5, 6)
	
	// These should use SIMD-friendly implementations internally
	addResult := c.EvalBinary(u, "+", v)
	subResult := c.EvalBinary(u, "-", v)
	mulResult := c.EvalBinary(u, "*", v)
	minResult := c.EvalBinary(u, "min", v)
	maxResult := c.EvalBinary(u, "max", v)
	
	// Verify results are vectors
	if _, ok := addResult.(*Vector); !ok {
		t.Errorf("add result is not a vector")
	}
	if _, ok := subResult.(*Vector); !ok {
		t.Errorf("sub result is not a vector")
	}
	if _, ok := mulResult.(*Vector); !ok {
		t.Errorf("mul result is not a vector")
	}
	if _, ok := minResult.(*Vector); !ok {
		t.Errorf("min result is not a vector")
	}
	if _, ok := maxResult.(*Vector); !ok {
		t.Errorf("max result is not a vector")
	}
	
	// Verify results are correct
	if !vectorsEqual(t, addResult.(*Vector), NewIntVector(5, 7, 9)) {
		t.Errorf("add integration result incorrect")
	}
	if !vectorsEqual(t, subResult.(*Vector), NewIntVector(-3, -3, -3)) {
		t.Errorf("sub integration result incorrect")
	}
	if !vectorsEqual(t, mulResult.(*Vector), NewIntVector(4, 10, 18)) {
		t.Errorf("mul integration result incorrect")
	}
	if !vectorsEqual(t, minResult.(*Vector), NewIntVector(1, 2, 3)) {
		t.Errorf("min integration result incorrect")
	}
	if !vectorsEqual(t, maxResult.(*Vector), NewIntVector(4, 5, 6)) {
		t.Errorf("max integration result incorrect")
	}
}
