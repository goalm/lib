package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func RecordToCsvString[T any](record T, suffix string) string {
	val := reflect.ValueOf(record)
	typ := reflect.TypeOf(record)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fields := suffix
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		res := f.String()
		switch f.Type().String() {
		case "int":
			res = strconv.Itoa(f.Interface().(int))
			if res == "" {
				res = "0"
			}
		case "float64":
			res = strconv.FormatFloat(f.Interface().(float64), 'f', 2, 64)
			if res == "" {
				res = "0.0"
			}

		case "string":
			if res == "" {
				res = "-"
			}
			res = `"` + res + `"`

		case "[]int":
			res = "["
			for _, v := range f.Interface().([]int) {
				res = res + strconv.Itoa(v) + " "
			}
			res = res + "]"

		case "[]float64":
			res = "["
			for _, v := range f.Interface().([]float64) {
				res = res + strconv.FormatFloat(v, 'f', 2, 64) + " "
			}
			res = res + "]"

		case "[]string":
			res = "["
			for _, v := range f.Interface().([]string) {
				res = res + `"` + v + `"` + " "
			}
			res = res + "]"
		}

		fields = fields + "," + res
	}
	return fields
}

func FieldsToCsvString[T any](a T, suffix string) string {
	val := reflect.ValueOf(a)
	typ := reflect.TypeOf(a)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	fields := suffix
	for i := 0; i < val.NumField(); i++ {
		fields = fields + "," + typ.Field(i).Name
	}
	return fields
}

func FindFieldByName[T any](s T, fieldName string) reflect.Value {
	val := reflect.ValueOf(s).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		fmt.Println("Field not found")
		return reflect.Value{}
	}
	return fieldVal
}

func IsFac(filePath string) bool {
	if strings.HasSuffix(filePath, ".fac") || strings.HasSuffix(filePath, ".FAC") {
		return true
	}
	return false
}

func computeLPSArray(pattern string) []int {
	var length = 0
	var i = 1
	var patternLength = len(pattern)

	var lps = make([]int, patternLength)

	lps[0] = 0

	for i = 1; i < patternLength; {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++

		} else {

			if length != 0 {
				length = lps[length-1]

			} else {
				lps[i] = length
				i++
			}
		}
	}
	return lps
}

func checkIfWholeWord(text string, startIndex int, endIndex int) bool {
	startIndex = startIndex - 1
	endIndex = endIndex + 1

	if (startIndex < 0 && endIndex >= len(text)) ||
		(startIndex < 0 && endIndex < len(text) && isNonWord(text[endIndex])) ||
		(startIndex >= 0 && endIndex >= len(text) && isNonWord(text[startIndex])) ||
		(startIndex >= 0 && endIndex < len(text) && isNonWord(text[startIndex]) && isNonWord(text[endIndex])) {
		return true
	}

	return false
}

func isNonWord(c byte) bool {
	return !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_'))
}

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
