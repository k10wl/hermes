package sqlite3

func convertToAnySlice[K ~[]E, E any](in K) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}
