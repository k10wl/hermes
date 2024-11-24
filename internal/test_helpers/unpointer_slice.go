package test_helpers

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

func ResetSliceTime[T timeResetter](arg []T) {
	for _, value := range arg {
		value.TimestampsToNilForTest__()
	}
}
