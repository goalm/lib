package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

// fastcache github.com/VictoriaMetrics/fastcache

type TblCache struct {
	Caches  *fastcache.Cache // cache for each table, with rowKey as keys
	SubKeys []string         // sub key for each record
}

func LoadFacToTblCache(filePath string, maxBytes int) *TblCache {
	start := time.Now()
	cache := fastcache.New(maxBytes)
	subKeys := make([]string, 0)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		noIdx := 0
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
					log.Fatalf("Table %v has no keys: %v", filePath, line)
				}
				str := strings.Split(line, ",")
				subKeys = str[noIdx:]
				cache.Set([]byte("subKeys"), []byte(strings.Join(subKeys, ",")))
				// process data
			} else if line[0] == '*' {
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				colVals := []byte(strings.Join(str[noIdx:], ","))
				cache.Set([]byte(key), colVals)
			}
		}
	} else {
		log.Fatalf("File %v is not a fac file", filePath)
	}
	fmt.Printf("loading data %s ... used %v\n", filePath, time.Since(start))

	return &TblCache{cache, subKeys}
}

func LoadFacHashToFastCache0814(filePath string, maxBytes int) ([]string, *fastcache.Cache) {
	start := time.Now()
	cache := fastcache.New(maxBytes)
	colKeys := make([]string, 0)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		noIdx := 0
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
					log.Fatalf("Table %v has no keys: %v", filePath, line)
				}
				str := strings.Split(line, ",")
				colKeys = str[noIdx:]
				cache.Set([]byte("colKeys"), []byte(strings.Join(colKeys, ",")))
			} else if line[0] == '*' {
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				colVals := []byte(strings.Join(str[noIdx:], ","))
				cache.Set([]byte(key), colVals)
			}
		}
	} else {
		log.Fatalf("File %v is not a fac file", filePath)
	}
	fmt.Printf("loading data %s ... used %v\n", filePath, time.Since(start))

	// size of cache

	return colKeys, cache
}

func LoadFacHashToFastCacheBk(filePath string, maxBytes int) *fastcache.Cache {
	start := time.Now()
	cache := fastcache.New(maxBytes)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		noIdx := 0
		colKeys := make([]string, 0)
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
					log.Fatalf("Table %v has no keys: %v", filePath, line)
				}
				str := strings.Split(line, ",")
				colKeys = str[noIdx:]
			} else if line[0] == '*' {
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}
				m := make(map[string]string)
				for i, v := range colKeys {
					m[v] = str[noIdx+i]
				}

				jsonStr, err := json.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}
				cache.Set([]byte(key), jsonStr)
			}
		}
	} else {
		log.Fatalf("File %v is not a fac file", filePath)
	}
	fmt.Printf("loading data %s ... used %v\n", filePath, time.Since(start))

	// size of cache

	return cache
}

func LoadFacToFastCache(filePath string, maxBytes int) *fastcache.Cache {
	start := time.Now()
	cache := fastcache.New(maxBytes)
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		noIdx := 0
		colKeys := make([]string, 0)
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
					log.Fatalf("Table %v has no keys: %v", filePath, line)
				}
				str := strings.Split(line, ",")
				colKeys = str[noIdx:]
			} else if line[0] == '*' {
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				for i, v := range colKeys {
					k := key + ":" + v
					cache.Set([]byte(k), []byte(str[noIdx+i]))
				}
			}
		}
	} else {
		log.Fatalf("File %v is not a fac file", filePath)
	}
	fmt.Printf("loading data %s ... used %v\n", filePath, time.Since(start))

	// size of cache

	return cache
}
