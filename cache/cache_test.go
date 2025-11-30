package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestCache_MaxTime(t *testing.T) {
	now := time.Now()
	later := now.Add(100 * time.Millisecond)
	laterNano := later.UnixNano()

	//fmt.Printf("Now: %v\nLater: %v\nLater Nano: %v\n", now, later, laterNano)

	seconds := laterNano / int64(time.Second)
	nanoseconds := laterNano % int64(time.Second)
	tm := time.Unix(seconds, nanoseconds)
	//fmt.Printf("Converted time: %v\n", tm)
	if !tm.Equal(later) {
		t.Errorf("Time conversion not equal: got = %v, want %v", tm, later)
	}

	tm2 := time.Unix(0, laterNano)
	//fmt.Printf("Converted time (alt method): %v\n", tm2)
	if !tm2.Equal(later) {
		t.Errorf("Time conversion (alt method) not equal: got = %v, want %v", tm2, later)
	}
}

func TestCache_Get(t *testing.T) {
	type args[K comparable] struct {
		k K
	}
	type testCase[K comparable, V any] struct {
		name  string
		c     *Cache[K, V]
		args  args[K]
		want  V
		want1 bool
	}

	// build our cache to test against
	cache := New[string, int]()
	cache.Set("always", 100)
	cache.SetTTL("expired", 200, 100*time.Millisecond)
	cache.SetTTL("valid", 300, 400*time.Millisecond)
	time.Sleep(150 * time.Millisecond) // make sure expired should be expired

	tests := []testCase[string, int]{
		{
			name:  "ok",
			c:     cache,
			args:  args[string]{k: "always"},
			want:  100,
			want1: true,
		},
		{
			name:  "expired",
			c:     cache,
			args:  args[string]{k: "expired"},
			want:  *new(int),
			want1: false,
		},
		{
			name:  "valid",
			c:     cache,
			args:  args[string]{k: "valid"},
			want:  300,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// re-test valid as it should now be invalid
	time.Sleep(275 * time.Millisecond)
	for _, tt := range tests {
		if tt.name == "valid" {
			t.Run("invalid", func(t *testing.T) {
				tt.want = *new(int)
				tt.want1 = false

				got, got1 := tt.c.Get(tt.args.k)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
				}
			})
		}
	}
}

func BenchmarkCache_Set(b *testing.B) {
	cache := New[string, []byte]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", []byte("value"))
	}
}

func BenchmarkCache_SetTTL(b *testing.B) {
	cache := New[string, []byte]()
	ttl := 60 * time.Second
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.SetTTL("key", []byte("value"), ttl)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cache := New[string, []byte]()
	cache.SetTTL("key", []byte("value"), 60*time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}
