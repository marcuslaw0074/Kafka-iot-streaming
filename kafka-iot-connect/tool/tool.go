package tool

import (
	"os"
	"reflect"
)

func WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func Contain[K comparable](l []K, k K) int {
	for ind, ele := range l {
		if reflect.DeepEqual(ele, k) {
			return ind
		}
	}
	return -1
}

func RemoveElementFromSlice[V any](slice []V, s int) []V {
	return append(slice[:s], slice[s+1:]...)
}