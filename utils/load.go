package utils

import (
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"io"
	"log"
	"os"
)

func LoadCsv[T any](fileName string, row T) []*T {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}
	var data []*T
	for {
		record := row
		if err := dec.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		data = append(data, &record)
	}

	return data
}

func ParseCsvMap(filename string) (map[string][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	keys := records[0]
	stringMap := make(map[string][]string, len(records)-1)
	for _, record := range records[1:] {
		for j, val := range record {
			if val != "" {
				stringMap[keys[j]] = append(stringMap[keys[j]], val)
			}

		}
	}

	return stringMap, nil
}

func MergeStringSlices(slices ...[]string) []string {
	uniqueStrings := make(map[string]bool)
	var mergedSlice []string
	for _, slice := range slices {
		for _, str := range slice {
			if !uniqueStrings[str] {
				uniqueStrings[str] = true
				mergedSlice = append(mergedSlice, str)
			}
		}
	}
	return mergedSlice
}

func StringExists(s []string, t string) bool {
	for _, str := range s {
		if str == t {
			return true
		}
	}
	return false
}
