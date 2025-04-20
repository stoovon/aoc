package maps

func Inverted[K, V comparable](m map[K]V) map[V]K {
	res := make(map[V]K, len(m))
	for k, v := range m {
		res[v] = k
	}
	return res
}

func MaxKey(m map[int]int) int {
	maxVal := -1
	for k := range m {
		if k > maxVal {
			maxVal = k
		}
	}
	return maxVal
}
