package test_helpers

func UnpointerSlice[T any](arg []*T) []T {
	res := make([]T, len(arg))
	for i, v := range arg {
		res[i] = *v
	}
	return res
}
