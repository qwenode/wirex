// Copyright 2018 The Wire Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wire

import (
    "context"
    "path/filepath"
    "runtime"
    "testing"
)

// BenchmarkGenerate benchmarks the standard Generate function.
func BenchmarkGenerate(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")
    opts := &GenerateOptions{}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := Generate(ctx, wd, nil, []string{"."}, opts)
        if len(errs) > 0 {
            b.Fatalf("Generate failed: %v", errs)
        }
    }
}

// BenchmarkGenerateOptimized benchmarks the optimized Generate function.
func BenchmarkGenerateOptimized(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")
    opts := &GenerateOptions{}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := GenerateOptimized(ctx, wd, nil, []string{"."}, opts)
        if len(errs) > 0 {
            b.Fatalf("GenerateOptimized failed: %v", errs)
        }
    }
}

// BenchmarkGenerateParallel benchmarks the parallel Generate function.
func BenchmarkGenerateParallel(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")
    opts := &GenerateOptions{}
    maxWorkers := runtime.GOMAXPROCS(0)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := GenerateParallel(ctx, wd, nil, []string{"."}, opts, maxWorkers)
        if len(errs) > 0 {
            b.Fatalf("GenerateParallel failed: %v", errs)
        }
    }
}

// BenchmarkProviderSetCache benchmarks the cache operations.
func BenchmarkProviderSetCache(b *testing.B) {
    cache := NewProviderSetCache()
    files := []string{"testdata/Chain/foo/wire.go"}

    // Create a dummy provider set for caching
    dummySet := &ProviderSet{
        PkgPath: "example.com/test",
        VarName: "TestSet",
    }

    b.Run("CacheSet", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            cache.CacheSet("example.com/test", "TestSet", dummySet, files)
        }
    })

    b.Run("GetCachedSet_Hit", func(b *testing.B) {
        cache.CacheSet("example.com/test", "TestSet", dummySet, files)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = cache.GetCachedSet("example.com/test", "TestSet", files)
        }
    })

    b.Run("GetCachedSetFast_Hit", func(b *testing.B) {
        cache.CacheSet("example.com/test", "TestSet", dummySet, files)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, _ = cache.GetCachedSetFast("example.com/test", "TestSet", files)
        }
    })

    b.Run("GetCachedSet_Miss", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _, _ = cache.GetCachedSet("example.com/nonexistent", "TestSet", files)
        }
    })
}

// BenchmarkLoad benchmarks the package loading function.
func BenchmarkLoad(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := load(ctx, wd, nil, "", []string{"."})
        if len(errs) > 0 {
            b.Fatalf("load failed: %v", errs)
        }
    }
}

// BenchmarkGenerateWithLazyLoad benchmarks the lazy loading Generate function.
func BenchmarkGenerateWithLazyLoad(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")
    opts := &GenerateOptions{}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := GenerateWithLazyLoad(ctx, wd, nil, []string{"."}, opts)
        if len(errs) > 0 {
            b.Fatalf("GenerateWithLazyLoad failed: %v", errs)
        }
    }
}

// BenchmarkGenerateParallelWithLazyLoad benchmarks parallel + lazy loading.
func BenchmarkGenerateParallelWithLazyLoad(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")
    opts := &GenerateOptions{}
    maxWorkers := runtime.GOMAXPROCS(0)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, errs := GenerateParallelWithLazyLoad(ctx, wd, nil, []string{"."}, opts, maxWorkers)
        if len(errs) > 0 {
            b.Fatalf("GenerateParallelWithLazyLoad failed: %v", errs)
        }
    }
}

// BenchmarkLazyLoadPackage benchmarks the lazy package loading.
func BenchmarkLazyLoadPackage(b *testing.B) {
    ctx := context.Background()
    wd := filepath.Join("testdata", "Chain", "foo")

    // First load the initial packages
    pkgs, errs := load(ctx, wd, nil, "", []string{"."})
    if len(errs) > 0 {
        b.Fatalf("load failed: %v", errs)
    }

    b.Run("WithLazyLoadEnabled", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            oc := newObjectCacheWithLazyLoad(pkgs, ctx, wd, nil)
            // Try to get a package that's already loaded (fast path)
            _, _ = oc.getPackage(pkgs[0].PkgPath)
        }
    })

    b.Run("WithoutLazyLoad", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            oc := newObjectCache(pkgs)
            _, _ = oc.getPackage(pkgs[0].PkgPath)
        }
    })
}
