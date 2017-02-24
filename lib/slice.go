package lib

//mapToSlice returns a slice of strings with the values of the map, the index is sorted
func StringMapToSlice(in map[string]string) []string {
	mapKeys := MapIndex(in)
	res := make([]string, len(in))
	for i, k := range mapKeys {
		res[i] = in[k]
	}
	return res
}
