package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type ProjResult struct {
	VarName string
	Res     []float64
}

func WriteCsv(start, end Date, res ProjResult) {
	f, err := os.Create("Results.csv")
	defer f.Close()

	f.WriteString("sep=,\n")

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()
	//header

	header := []string{"VAR_NAME"}
	n := Dur(start, end)
	for i := 0; i <= n; i++ {
		calDate := strconv.Itoa(start.CalendarDate(i).DateValue())
		header = append(header, calDate)
	}

	if err := w.Write(header); err != nil {
		log.Fatalln("error writing header to file", err)
	}

	record := []string{res.VarName}
	for _, value := range res.Res {
		record = append(record, strconv.FormatFloat(value, 'f', 2, 64))
	}
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}
