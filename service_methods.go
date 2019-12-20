package testtask1

func InArray(val string, array []string) bool {
	for _, row := range array {
		if val == row {
			return true
		}
	}

	return false
}
