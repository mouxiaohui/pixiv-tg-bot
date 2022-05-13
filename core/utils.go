package core

func removeArrVal[T string | int](data []T, target T) []T {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			data = append(data[:i], data[i+1:]...)
		}
	}

	return data
}
