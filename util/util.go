package util

import (
	"reflect"
	"engo.io/engo"
	"engo.io/engo/common"
)

//panic if s is not a slice
func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func LoadTmxMap(filePath string) (*common.Level, error) {
	err := engo.Files.Load(filePath)
	if err != nil {
		return nil, err
	}
	resource, err := engo.Files.Resource(filePath)
	if err != nil {
		return nil, err
	}
	// We need the return above, we would like to avoid cast error
	level := resource.(common.TMXResource).Level
	return level, nil
}