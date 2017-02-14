package gogen

import "sort"

//mapToSlice returns a slice of strings with the values of the map, the index is sorted
func stringMapToSlice(in map[string]string) []string {
	mapKeys := mapIndex(in)
	res := make([]string, len(in))
	for i, k := range mapKeys {
		res[i] = in[k]
	}
	return res
}

//mapIndex retuns a slice of strings with all the map key indexes sorted
func mapIndex(in map[string]string) []string {

	index := make([]string, 0, len(in))
	for key, _ := range in {
		index = append(index, key)
	}
	sort.Strings(index)
	return index
}
