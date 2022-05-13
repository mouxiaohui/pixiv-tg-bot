package core

func removeArrVal[T string | int](data []T, target T) []T {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			data = append(data[:i], data[i+1:]...)
		}
	}

	return data
}

func arrToString(arr []string, separator string) string {
	var res string
	for i, a := range arr {
		res += a
		if i != len(arr)-1 {
			res += separator
		}
	}

	return res
}
