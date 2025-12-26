# SIMD Implementation in Ivy

This implementation uses Go's experimental SIMD API (https://github.com/golang/go/issues/73787) to accelerate vector operations on Int values.

## Requirements

- **Go 1.25 or later** (currently in development on master branch)
- Build with `GOEXPERIMENT=simd`

**Note**: The SIMD experiment is not yet available in released Go versions. To use it, you need to build Go from source from the master branch.

## Building Go from Source

To build Go with SIMD support:

```bash
# Clone Go repository
git clone https://go.googlesource.com/go
cd go
git checkout master

# Build Go
cd src
./make.bash

# Use the built Go
export PATH=$(pwd)/../bin:$PATH
go version  # Should show go1.25 or later
```

## Building with SIMD Support

To build Ivy with SIMD support:

```bash
GOEXPERIMENT=simd go build
```

To run tests with SIMD:

```bash
GOEXPERIMENT=simd go test ./...
```

To run benchmarks:

```bash
GOEXPERIMENT=simd go test -bench=VectorAdd ./value -benchmem
GOEXPERIMENT=simd go test -bench=. ./value -benchmem
```

## Implementation Details

### SIMD Operations

The following vector operations have SIMD implementations:

- **Addition (+)**: Uses `Int64x4.Add()` (AVX2) or `Int64x8.Add()` (AVX512)
- **Subtraction (-)**: Uses `Int64x4.Sub()` (AVX2) or `Int64x8.Sub()` (AVX512)
- **Multiplication (*)**: Uses `Int64x8.Mul()` (AVX512DQ required)
- **Minimum (min)**: Uses `Int64x8.Min()` (AVX512F)
- **Maximum (max)**: Uses `Int64x8.Max()` (AVX512F)

### CPU Feature Detection

The implementation checks for CPU features at runtime:
- `archsimd.X86.HasAVX2()` - 256-bit vectors (Int64x4)
- `archsimd.X86.HasAVX512F()` - 512-bit vectors (Int64x8)
- `archsimd.X86.HasAVX512DQ()` - Required for int64 multiplication

### Vector Processing

- **AVX2**: Processes 4 int64 values per iteration
- **AVX512**: Processes 8 int64 values per iteration
- Remaining elements are processed with scalar operations

### Automatic Fallback

SIMD operations automatically fall back to generic implementations when:
1. SIMD experiment is not enabled (build without `GOEXPERIMENT=simd`)
2. Vectors contain non-Int types (BigInt, BigRat, BigFloat, Complex)
3. CPU doesn't support required SIMD instructions
4. Vector length is too small to benefit from SIMD

## Architecture Support

Currently, SIMD operations are implemented for:
- **AMD64** (x86-64): SSE, AVX2, AVX512

The `simd/archsimd` package provides architecture-specific implementations.

## Performance Expectations

Expected speedup depends on:
- Vector size (larger = better)
- CPU architecture and SIMD support level
- Memory bandwidth
- Data alignment

Typical improvements on AVX2/AVX512 CPUs:
- Small vectors (< 100 elements): Minimal or negative due to setup overhead
- Medium vectors (100-10000): 1.5-3x speedup
- Large vectors (> 10000): 2-4x speedup

## Code Structure

```
value/
  simd_ops.go          - SIMD implementations (GOEXPERIMENT=simd)
  simd_ops_nosimd.go   - Stub implementations (no GOEXPERIMENT)
  simd_bench_test.go   - Benchmarks
  eval.go              - Integration with binaryVectorOp
```

## Example Usage

From Ivy's perspective, SIMD is completely transparent:

```ivy
# These automatically use SIMD when available
(1 2 3 4 5) + (10 20 30 40 50)
(1 2 3) * (4 5 6)
(1 5 3) min (4 2 6)
```

## Verifying SIMD Usage

To verify SIMD instructions are being generated:

```bash
# Build with assembly output
GOEXPERIMENT=simd go build -gcflags="-S" ./value 2>&1 | grep -A 20 vectorAddSIMD

# Look for SIMD instructions:
# - VPADDD, VPADDQ (AVX2 add)
# - VPSUBQ (AVX2 subtract)
# - VPMULLQ (AVX512DQ multiply)
# - VPMINSQ, VPMAXSQ (AVX512 min/max)
```

## References

- [Go SIMD Experiment Proposal](https://github.com/golang/go/issues/73787)
- [simd/archsimd Package Documentation](https://go.googlesource.com/go/+/refs/heads/master/src/simd/archsimd/doc.go)
- [Intel Intrinsics Guide](https://www.intel.com/content/www/us/en/docs/intrinsics-guide/index.html)

## Limitations

1. **Experimental API**: The SIMD API is experimental and subject to change
2. **Architecture-specific**: Only AMD64 currently supported
3. **Type restrictions**: Only works with Int vectors (not BigInt/BigRat/BigFloat)
4. **CPU requirements**: Requires AVX2 or AVX512 for meaningful speedups
5. **Overflow handling**: Still requires `maybeBig()` checks, adding overhead

## Future Work

Potential improvements:
1. Support for 32-bit integers (Int32)
2. ARM NEON support
3. More operations (bitwise, shifts)
4. Better batching to amortize setup costs
5. Alignment hints for better performance
