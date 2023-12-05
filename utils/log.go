package utils

import (
	"log"
	"os"
)

func SetLogOutput() {
	file := "./" + "log" + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}
