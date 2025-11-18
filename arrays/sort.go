package arrays

// SubarraySort returns starting and ending indices of arr
// in order to make all the array sorted
func SubarraySort(arr []int) []int {
	indices := []int{-1, -1}
	for i, mx := 1, arr[0]; i < len(arr); i++ {
		if mx > arr[i] {
			for k := 0; k < i; k++ {
				if arr[k] > arr[i] {
					if indices[0] == -1 || indices[0] > k {
						indices[0] = k
					}
					indices[1] = i
				}
			}
		} else {
			mx = arr[i]
		}
	}
	return indices
}
