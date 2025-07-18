package util

func ValueToPointer[T any](value T) *T {
	return &value
}
