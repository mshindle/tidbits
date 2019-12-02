package toy

func f(left, right chan int) {
	left <- 1 + <-right
}

// Whisper creates N go routines which takes an integer from the right hand
// neighbor, adds 1, and passes it to the left hand neighbor.
func Whisper(n int) int {
	var right chan int

	leftmost := make(chan int)
	left := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go f(left, right)
		left = right
	}
	go func(c chan int) { c <- 0 }(right)
	return <-leftmost
}
