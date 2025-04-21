package maths

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(nums ...int) int {
	if len(nums) == 0 {
		panic("Max requires at least one argument")
	}
	maximum := nums[0]
	for _, n := range nums[1:] {
		if n > maximum {
			maximum = n
		}
	}
	return maximum
}

func Min(nums ...int) int {
	if len(nums) == 0 {
		panic("Min requires at least one argument")
	}
	minimum := nums[0]
	for _, n := range nums[1:] {
		if n < minimum {
			minimum = n
		}
	}
	return minimum
}

// Sign returns the sign of an integer (-1, 0, or 1).
func Sign(i int) int {
	if i < 0 {
		return -1
	} else if i > 0 {
		return 1
	}
	return 0
}
