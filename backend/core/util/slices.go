package util

// SliceLastNElements returns the last n elements of a slice.
// If the length of the slice is less than n, it returns all elements.
func SliceLastNElements[T any](slice []T, n int) []T {
	if len(slice) < n {
		return slice
	}
	return slice[len(slice)-n:]
}
