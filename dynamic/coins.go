package dynamic

import (
	"math"
)

// Coins finds the minimum number of coins to add up to value
func Coins(n int, denoms []int) int {
	// min as calculated for each increment
	numCoins := make([]int, n+1)
	for i := 1; i <= n; i++ {
		numCoins[i] = math.MaxInt32
	}
	numCoins[0] = 0

	for i := range numCoins {
		for _, denom := range denoms {
			if denom <= i && numCoins[i-denom]+1 < numCoins[i] {
				numCoins[i] = numCoins[i-denom] + 1
			}
		}
	}

	if numCoins[n] == math.MaxInt32 {
		return -1
	}
	return numCoins[n]
}
