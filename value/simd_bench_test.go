// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"testing"
	
	"robpike.io/ivy/config"
)

// Simple mock context for benchmarking that doesn't require exec package
type benchContext struct {
	config *config.Config
}

func (b benchContext) Config() *config.Config {
	return b.config
}

func (b benchContext) Local(i int) *Var {
	return nil
}

func (b benchContext) Global(name string) *Var {
	return nil
}

func (b benchContext) AssignGlobal(name string, value Value) {
}

func (b benchContext) Eval(exprs []Expr) []Value {
	return nil
}

func (b benchContext) EvalUnary(op string, right Value) Value {
	return UnaryOps[op].EvalUnary(b, right)
}

func (b benchContext) EvalBinary(left Value, op string, right Value) Value {
	return BinaryOps[op].EvalBinary(b, left, right)
}

func (b benchContext) UserDefined(op string, isBinary bool) bool {
	return false
}

func (b benchContext) TraceIndent() string {
	return ""
}

// Helper to create Int vectors
func makeIntVector(size int) *Vector {
	edit := newVectorEditor(size, nil)
	for i := 0; i < size; i++ {
		edit.Set(i, Int(i))
	}
	return edit.Publish()
}

// Helper to create Int vectors with specific values
func makeIntVectorValue(size int, val Int) *Vector {
	edit := newVectorEditor(size, nil)
	for i := 0; i < size; i++ {
		edit.Set(i, val)
	}
	return edit.Publish()
}

func getBenchContext() Context {
	conf := &config.Config{}
	// Config will initialize itself on first use
	return benchContext{config: conf}
}

// Benchmarks for vector addition

func BenchmarkVectorAddOriginal_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAddSIMD_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorAddInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD add returned nil")
		}
	}
}

func BenchmarkVectorAddOriginal_Medium(b *testing.B) {
	u := makeIntVector(1000)
	v := makeIntVector(1000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAddSIMD_Medium(b *testing.B) {
	u := makeIntVector(1000)
	v := makeIntVector(1000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorAddInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD add returned nil")
		}
	}
}

func BenchmarkVectorAddOriginal_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAddSIMD_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorAddInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD add returned nil")
		}
	}
}

// Benchmarks for vector subtraction

func BenchmarkVectorSubOriginal_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "-", v)
	}
}

func BenchmarkVectorSubSIMD_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorSubInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD sub returned nil")
		}
	}
}

func BenchmarkVectorSubOriginal_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "-", v)
	}
}

func BenchmarkVectorSubSIMD_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorSubInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD sub returned nil")
		}
	}
}

// Benchmarks for vector multiplication

func BenchmarkVectorMulOriginal_Small(b *testing.B) {
	u := makeIntVectorValue(100, 2)
	v := makeIntVectorValue(100, 3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "*", v)
	}
}

func BenchmarkVectorMulSIMD_Small(b *testing.B) {
	u := makeIntVectorValue(100, 2)
	v := makeIntVectorValue(100, 3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMulInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD mul returned nil")
		}
	}
}

func BenchmarkVectorMulOriginal_Large(b *testing.B) {
	u := makeIntVectorValue(10000, 2)
	v := makeIntVectorValue(10000, 3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "*", v)
	}
}

func BenchmarkVectorMulSIMD_Large(b *testing.B) {
	u := makeIntVectorValue(10000, 2)
	v := makeIntVectorValue(10000, 3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMulInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD mul returned nil")
		}
	}
}

// Benchmarks for vector min

func BenchmarkVectorMinOriginal_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "min", v)
	}
}

func BenchmarkVectorMinSIMD_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMinInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD min returned nil")
		}
	}
}

func BenchmarkVectorMinOriginal_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "min", v)
	}
}

func BenchmarkVectorMinSIMD_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMinInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD min returned nil")
		}
	}
}

// Benchmarks for vector max

func BenchmarkVectorMaxOriginal_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "max", v)
	}
}

func BenchmarkVectorMaxSIMD_Small(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMaxInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD max returned nil")
		}
	}
}

func BenchmarkVectorMaxOriginal_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "max", v)
	}
}

func BenchmarkVectorMaxSIMD_Large(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorMaxInt(c, u, v)
		if result == nil {
			b.Fatal("SIMD max returned nil")
		}
	}
}

// Benchmarks for scalar operations

func BenchmarkVectorScalarAddSIMD_Small(b *testing.B) {
	v := makeIntVector(100)
	scalar := Int(42)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorScalarAddInt(c, scalar, v)
		if result == nil {
			b.Fatal("SIMD scalar add returned nil")
		}
	}
}

func BenchmarkVectorScalarAddSIMD_Large(b *testing.B) {
	v := makeIntVector(10000)
	scalar := Int(42)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorScalarAddInt(c, scalar, v)
		if result == nil {
			b.Fatal("SIMD scalar add returned nil")
		}
	}
}

func BenchmarkVectorScalarMulSIMD_Small(b *testing.B) {
	v := makeIntVector(100)
	scalar := Int(3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorScalarMulInt(c, scalar, v)
		if result == nil {
			b.Fatal("SIMD scalar mul returned nil")
		}
	}
}

func BenchmarkVectorScalarMulSIMD_Large(b *testing.B) {
	v := makeIntVector(10000)
	scalar := Int(3)
	c := getBenchContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := vectorScalarMulInt(c, scalar, v)
		if result == nil {
			b.Fatal("SIMD scalar mul returned nil")
		}
	}
}

