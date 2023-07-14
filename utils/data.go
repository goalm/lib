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
