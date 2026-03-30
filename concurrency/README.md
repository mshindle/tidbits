# concurrency

Sample functions and tests that highlight common concurrency mistakes and provide patterns for solving them.

## Leaking goroutines

**Problem:** Long-lived goroutines wait on channels that never close.

**Example:**
```go
package main

import (
 "time"
)
func watch(ch <-chan struct{}) {
 for {
  select {
  case <-ch:
   // intended stop
  }
 }
}

func main() {
 for i := 0; i < 100; i++ {
  c := make(chan struct{})
  go watch(c) // never closed
 }
 time.Sleep(1 * time.Second)
}
```

**Solution:** Use context cancellation to signal goroutines to exit when they are no longer needed. See leaking.go for an example.


