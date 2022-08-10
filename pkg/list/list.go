package list

func Contains[T comparable](arr []T, v T) bool {
	for _, x := range arr {
		if x == v {
			return true
		}
	}
	return false
}
