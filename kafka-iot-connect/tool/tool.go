package tool

import (
	"encoding/json"
	"io"
	"log"
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

func ParseJsonFile[V any](file string) (res V, err error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return
	} else {
		defer jsonFile.Close()
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			return res, err
		} else {
			return res, json.Unmarshal(byteValue, &res)
		}
	}
}

func Filter[S ~[]E, E any](s S, f func(E, int) bool) S {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	ls := make(S, 0)
	for ind, ele := range s {
		if f(ele, ind) {
			ls = append(ls, ele)
		}
	}
	return ls
}

func Map[S ~[]E, E, K any](s S, f func(E, int) K) []K {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	ls := make([]K, 0)
	for ind, ele := range s {
		ls = append(ls, f(ele, ind))
	}
	return ls
}

func UniqueElements[S ~[]T, T comparable](s S) S {
	unique := make(map[T]bool, len(s))
	us := make([]T, len(unique))
	for _, elem := range s {
		if !unique[elem] {
			us = append(us, elem)
			unique[elem] = true
		}
	}
	return us
}

func ReverseKeyValue[V comparable, K comparable](m map[K]V) map[V]K {
	mm := map[V]K{}
	for k, v := range m {
		mm[v] = k
	}
	return mm
}

func RemoveDuplicate[S ~[]T, T comparable](sliceList S) S {
    allKeys := make(map[T]bool)
    list := []T{}
    for _, item := range sliceList {
        if _, value := allKeys[item]; !value {
            allKeys[item] = true
            list = append(list, item)
        }
    }
    return list
}