package lib

import "sort"

//MapIndex retuns a slice of strings with all the map keys indexes sorted
func MapIndex(in map[string]string) []string {
	index := make([]string, 0, len(in))
	for key, _ := range in {
		index = append(index, key)
	}
	sort.Strings(index)
	return index
}
