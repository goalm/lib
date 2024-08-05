package utils

import (
	"fmt"
	"reflect"
	"regexp"
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
			res = ""
			cap := val.Field(i).Cap()
			s := f.Interface().([]int)
			for j := 0; j < cap-1; j++ {
				res = res + strconv.Itoa(s[j]) + ","
			}
			res = res + strconv.Itoa(s[cap-1])

		case "[]float64":
			res = ""
			cap := val.Field(i).Cap()
			if cap == 0 {
				res = "0.0"
			} else {
				s := f.Interface().([]float64)
				for j := 0; j < cap-1; j++ {
					res = res + strconv.FormatFloat(s[j], 'f', 2, 64) + ","
				}
				//todo: remove max(0, ...)
				res = res + strconv.FormatFloat(s[cap-1], 'f', 2, 64)
			}

		case "[]string":
			res = ""
			cap := val.Field(i).Cap()
			s := f.Interface().([]string)
			for j := 0; j < cap-1; j++ {
				res = res + s[j] + ","
			}
			res = res + s[cap-1]
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
		if typ.Field(i).Type.Kind() != reflect.Slice {
			fields = fields + "," + typ.Field(i).Name
		} else {
			for j := 0; j < val.Field(i).Cap(); j++ {
				fields = fields + "," + typ.Field(i).Name + "(" + strconv.Itoa(j+1) + ")"
			}
		}
	}
	return fields
}

func FindFieldByName[T any](s T, fieldName string) (reflect.Value, bool) {
	val := reflect.ValueOf(s).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		fmt.Println("Field not found")
		return reflect.Value{}, false
	}
	return fieldVal, true
}

func IsFac(filePath string) bool {
	if strings.HasSuffix(filePath, ".fac") || strings.HasSuffix(filePath, ".FAC") {
		return true
	}
	return false
}

func FilePathToName(str string) (string, error) {
	re := regexp.MustCompile(`[\/|\\]*(\w*)\.\w*$`)
	match := re.FindStringSubmatch(str)
	if len(match) > 0 {
		return match[1], nil
	}
	return "", fmt.Errorf("No match found")
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
