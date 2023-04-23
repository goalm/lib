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
