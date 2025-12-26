// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"testing"
	
	"robpike.io/ivy/config"
)

// Helper functions for benchmarking

func makeIntVector(size int) *Vector {
	edit := newVectorEditor(size, nil)
	for i := 0; i < size; i++ {
		edit.Set(i, Int(i))
	}
	return edit.Publish()
}

func makeIntVectorValue(size int, val Int) *Vector {
	edit := newVectorEditor(size, nil)
	for i := 0; i < size; i++ {
		edit.Set(i, val)
	}
	return edit.Publish()
}

// Benchmarks for vector addition with SIMD

func BenchmarkVectorAdd_Small_100(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAdd_Medium_1000(b *testing.B) {
	u := makeIntVector(1000)
	v := makeIntVector(1000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAdd_Large_10000(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

func BenchmarkVectorAdd_XLarge_100000(b *testing.B) {
	u := makeIntVector(100000)
	v := makeIntVector(100000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "+", v)
	}
}

// Benchmarks for vector subtraction

func BenchmarkVectorSub_Small_100(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "-", v)
	}
}

func BenchmarkVectorSub_Large_10000(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "-", v)
	}
}

// Benchmarks for vector multiplication

func BenchmarkVectorMul_Small_100(b *testing.B) {
	u := makeIntVectorValue(100, 2)
	v := makeIntVectorValue(100, 3)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "*", v)
	}
}

func BenchmarkVectorMul_Large_10000(b *testing.B) {
	u := makeIntVectorValue(10000, 2)
	v := makeIntVectorValue(10000, 3)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "*", v)
	}
}

// Benchmarks for vector min

func BenchmarkVectorMin_Small_100(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "min", v)
	}
}

func BenchmarkVectorMin_Large_10000(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "min", v)
	}
}

// Benchmarks for vector max

func BenchmarkVectorMax_Small_100(b *testing.B) {
	u := makeIntVector(100)
	v := makeIntVector(100)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "max", v)
	}
}

func BenchmarkVectorMax_Large_10000(b *testing.B) {
	u := makeIntVector(10000)
	v := makeIntVector(10000)
	c := testContextForBench()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.EvalBinary(u, "max", v)
	}
}

// Test context helper
func testContextForBench() Context {
	return &benchContextImpl{}
}

type benchContextImpl struct{}

func (b *benchContextImpl) Config() *config.Config {
	return &config.Config{}
}

func (b *benchContextImpl) Local(i int) *Var {
	return nil
}

func (b *benchContextImpl) Global(name string) *Var {
	return nil
}

func (b *benchContextImpl) AssignGlobal(name string, value Value) {
}

func (b *benchContextImpl) Eval(exprs []Expr) []Value {
	return nil
}

func (b *benchContextImpl) EvalUnary(op string, right Value) Value {
	return UnaryOps[op].EvalUnary(b, right)
}

func (b *benchContextImpl) EvalBinary(left Value, op string, right Value) Value {
	return BinaryOps[op].EvalBinary(b, left, right)
}

func (b *benchContextImpl) UserDefined(op string, isBinary bool) bool {
	return false
}

func (b *benchContextImpl) TraceIndent() string {
	return ""
}
