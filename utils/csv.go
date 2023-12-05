package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func ReadProdSpecs(filePath string) ([]string, map[string][]string) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	header, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	dataMap := make(map[string][]string)
	pNames := make([]string, 0)
	for i, h := range header {
		if i > 0 {
			dataMap[h] = make([]string, 0)
			pNames = append(pNames, h)
		}

	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		ind := record[0]
		for i, value := range record {
			if i > 0 {
				if value == "*" {
					value = ind
				} else {
					continue
				}
				header := header[i]
				dataMap[header] = append(dataMap[header], value)
			}
		}
	}
	return pNames, dataMap
}
