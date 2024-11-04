package utils

func GetIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}

	return -1
}
