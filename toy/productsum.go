package toy

type SpecialArray []interface{}

// Write a function that takes in a "special" array and returns its product sum. A "special" array is a non-empty
// array that contains either integers or other "special" arrays. The product sum of a "special" array is the sum of
// its elements, where "special" arrays inside it should be summed themselves and then multiplied by their level of
// depth. For example, the product sum of [x, y] is x + y; the product sum of [x, [y, z]] is x + 2y + 2z.
func ProductSum(array []interface{}) int {
	return addSpecial(array, 1)
}

func addSpecial(array SpecialArray, level int) int {
	var sum int
	for _, i := range array {
		switch e := i.(type) {
		case int:
			sum += e
		case SpecialArray:
			sum += addSpecial(e, level+1)
		}
	}
	return level * sum
}
