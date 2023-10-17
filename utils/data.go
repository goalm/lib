package utils

import (
	"fmt"
	"reflect"
)

func RemoveDup(this []*float64) []*float64 {
	tgt := this[:1]
	for _, v := range this {
		tgt = append(tgt, v)
	}
	return tgt
}

func MapUniqValues[T comparable](m map[string]T) []T {
	var uniqueValues []T
	uniqueMap := make(map[T]bool)
	for _, value := range m {
		if _, ok := uniqueMap[value]; !ok {
			uniqueValues = append(uniqueValues, value)
			uniqueMap[value] = true
		}
	}

	return uniqueValues
}

func MapContainsValue[K, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}

	return false
}

func PivotFloat64(data any, tags []string, fieldsToSum []string) map[string]map[string]float64 {
	pivot := make(map[string]map[string]float64)
	listValue := reflect.ValueOf(data)

	for i := 0; i < listValue.Len(); i++ {
		var tagKey string
		for _, tag := range tags {
			tagValue := reflect.Indirect(listValue.Index(i)).FieldByName(tag)
			tagKey += fmt.Sprintf("%v", tagValue.Interface())
		}

		if _, ok := pivot[tagKey]; !ok {
			pivot[tagKey] = make(map[string]float64)
		}

		for _, sum := range fieldsToSum {
			sumValue := reflect.Indirect(listValue.Index(i)).FieldByName(sum)
			pivot[tagKey][sum] += sumValue.Float()
		}

		pivot[tagKey]["count"] = pivot[tagKey]["count"] + 1
	}

	return pivot
}
