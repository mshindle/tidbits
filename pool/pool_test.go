package pool

import (
	"sync"
	"testing"
)

func BenchmarkSlicePool(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]int, 0, 0)
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := pool.Get().([]int)
		s = s[0:0]
		s = append(s, 123)
		pool.Put(s)
	}
}

func BenchmarkSlicePoolPtr(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			s := make([]int, 0, 0)
			return &s
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := pool.Get().(*[]int)
		s := *ptr
		s = s[0:0]
		s = append(s, 123)
		*ptr = s
		pool.Put(ptr)
	}
}
