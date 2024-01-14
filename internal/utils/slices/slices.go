package slices

// ProcessSubsetDifference runs a process function on the elements of a that are not in b.
func ProcessSubsetDifference[T comparable, V any](b []T, a []T, processFunc func(T) V) []V {
	m := make(map[T]bool)

	for _, v := range b {
		m[v] = true
	}

	var diff []V

	for _, v := range a {
		if _, ok := m[v]; !ok {
			diff = append(diff, processFunc(v))
		}
	}

	return diff
}
