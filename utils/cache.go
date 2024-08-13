package utils

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/allegro/bigcache/v3"
)

// fastcache github.com/VictoriaMetrics/fastcache
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

func LoadFacToBigCache(filePath string) *bigcache.BigCache {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	cache.Set("my-unique-key", []byte(filePath))
	value, _ := cache.Get("my-unique-key")
	fmt.Printf("value is %s\n", string(value))
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
					cache.Set(k, []byte(str[noIdx+i]))
				}
			}
		}
	} else {
		log.Fatalf("File %v is not a fac file", filePath)
	}

	return cache
}
