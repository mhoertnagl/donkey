package utils

func Map[S, T any](values []S, f func(S) T) []T {
	res := make([]T, len(values))
	for i, v := range values {
		res[i] = f(v)
	}
	return res
}
