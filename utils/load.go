package utils

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jszwec/csvutil"
)

func LoadCsvToEnum(filePath string) *Enum {
	m := NewEnum()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, record := range records {
		if i == 0 {
			m.SetNames(record[0], record[1])

		} else {
			k, err := strconv.Atoi(record[0])
			if err != nil {
				log.Fatal(err)
			}
			m.Add(k, record[1])
		}
	}
	return m
}

func LoadFacToMap(filePath string) map[string]string {
	m := make(map[string]string)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		noIdx := 0
		hashKeys := make([]string, 0)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			// end of file
			if line == "\xA0" || line == "" {
				break
			}
			// dump descriptions
			if line[0] != '!' && line[0] != '*' {
				continue
			}

			line = strings.ReplaceAll(line, "\"", "")
			// process header
			if line[0] == '!' {
				noIdx, err = strconv.Atoi(line[1:2])
				if err != nil {
					log.Fatal(err)
				}

				if noIdx < 1 {
					log.Fatal("Table has no keys")
				}

				str := strings.Split(line, ",")
				hashKeys = str[noIdx:]

			} else if line[0] == '*' {
				// Records
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				for i, v := range hashKeys {
					m[key+":"+v] = str[noIdx+i]
				}
			}
		}

	} else {
		log.Fatal("File is not a .fac file")
	}

	return m
}

func LoadFacToHashMap(filePath string) map[string]map[string]string {
	m := make(map[string]map[string]string)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		noIdx := 0
		hashKeys := make([]string, 0)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ReplaceAll(line, "\"", "")
			// dump descriptions
			if line[0] != '!' && line[0] != '*' {
				continue
			}
			// process header
			if line[0] == '!' {
				noIdx, err = strconv.Atoi(line[1:2])
				if err != nil {
					log.Fatal(err)
				}

				if noIdx < 1 {
					log.Fatal("Table has no keys")
				}

				str := strings.Split(line, ",")
				hashKeys = str[noIdx:]

			} else if line[0] == '*' {
				// Records
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				h := make(map[string]string)
				for i, v := range hashKeys {
					h[v] = str[noIdx+i]
				}
				m[key] = h
			}
		}

	} else {
		log.Fatal("File is not a .fac file")
	}

	return m
}

func LoadPropMpToChn[T any](fileName string, dataStruct T, dataChn chan *T) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(str, "VARIABLE_TYPES") {
			break
		}
	}

	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record := dataStruct
		if err := dec.Decode(&record); err == io.EOF {
			close(dataChn)
			break
		} else if err != nil {
			log.Println("Error reading " + fileName + ": " + err.Error())
			close(dataChn)
			break
		}
		dataChn <- &record
	}
}

func LoadPropMpToStruct[T any](fileName string, dataStruct T) []*T {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	skipped := false

	if !skipped {
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasPrefix(str, "VARIABLE_TYPES") {
				break
			}
		}
		skipped = true
	}

	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}
	var data []*T
	for {
		record := dataStruct
		if err := dec.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error reading " + fileName + ": " + err.Error())
			break
		}

		data = append(data, &record)
	}
	return data
}

func LoadFileToStruct[T any](fileName string, skipLines int, dataStruct T) []*T {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for i := 0; i < skipLines; i++ {
		_, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	}

	csvReader := csv.NewReader(reader)
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		log.Fatal(err)
	}
	var data []*T
	for {
		record := dataStruct
		if err := dec.Decode(&record); err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error reading " + fileName + ": " + err.Error())
			break
		}

		data = append(data, &record)
	}
	return data
}

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
			log.Println("Error reading " + fileName + ": " + err.Error())
			break
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
