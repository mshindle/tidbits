package arrays

import "sort"

func intArraySort(arr [][]int) func(i, j int) bool {
	return func(i, j int) bool {
		var indx int
		if arr[i][indx] == arr[j][indx] {
			indx = 1
		}
		return arr[i][indx] < arr[j][indx]
	}
}

// TwoNumberSum returns an array with all the digit pairs which add to target
func TwoNumberSum(array []int, target int) [][]int {
	nums := len(array)
	results := make([][]int, 0, nums)
	for x := 0; x < nums; x++ {
		for y := x + 1; y < nums; y++ {
			if array[x]+array[y] == target {
				pair := []int{array[x], array[y]}
				sort.Ints(pair)
				results = append(results, pair)
			}
		}
	}
	sort.Slice(results, intArraySort(results))
	return results
}

// TwoNumberSum returns an array with all the digit triplets which add to target
func ThreeNumberSum(array []int, target int) [][]int {
	nums := len(array)
	results := make([][]int, 0, nums)
	for x := 0; x < nums; x++ {
		for y := x + 1; y < nums; y++ {
			need := target - (array[x] + array[y])
			for z := y + 1; z < nums; z++ {
				if array[z] == need {
					triplet := []int{array[x], array[y], array[z]}
					sort.Ints(triplet)
					results = append(results, triplet)
				}
			}
		}
	}
	sort.Slice(results, intArraySort(results))
	return results
}
