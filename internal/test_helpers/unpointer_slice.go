package test_helpers

import "slices"

func UnpointerSlice[T any](arg []*T) []T {
	res := make([]T, len(arg))
	for i, v := range arg {
		res[i] = *v
	}
	return res
}

type timeResetter interface {
	TimestampsToNilForTest__()
}

func ResetSliceTime[T timeResetter](s []T) []T {
	res := slices.Clone(s)
	for _, value := range res {
		value.TimestampsToNilForTest__()
	}
	return res
}
