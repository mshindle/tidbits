package random

import (
	"math/rand"
)

// pseudoFlip returns a pesudo-random true / false.
func pseudoFlip() bool {
	return rand.Int()%2 == 0
}
