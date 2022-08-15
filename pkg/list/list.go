package list

func Contains[T comparable](arr []T, v T) bool {
	for _, x := range arr {
		if x == v {
			return true
		}
	}
	return false
}

// TakeN takes  at most N values from the given offset
func TakeN[T any](arr []T, n int, offset int) []T {
	if len(arr) <= offset {
		return nil
	}
	if len(arr)-offset < n {
		n = len(arr) - offset
	}
	return arr[offset : offset+n]
}
