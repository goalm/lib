package utils

import (
	"fmt"
	"log"
	"os"
)

func GetFile(fileName string, paths []string) (fileLoc string) {
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Printf("File %s not found in any of the paths\n", fileName)
	return
}

// todo: remove the below functions
func FindMp(fileName string) (fileLoc string) {
	for _, path := range MpLocs {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Println("Model Point " + fileName + " not found in any of the paths")
	return
}

func FindFile(tbl string) (fileLoc string) {
	fileName := Conf.GetString("fileNames." + tbl)
	for _, path := range TableLocs {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory " + path)
			continue
		}
		for _, file := range files {
			if file.Name() == fileName {
				fileLoc = path + "/" + fileName
				log.Println(fileName+" found at: ", fileLoc)
				return
			}
		}
	}
	fmt.Println("Table " + fileName + " not found in any of the paths")
	return
}
