package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var once sync.Once

func GetRedisClient() *redis.Client {
	once.Do(func() {
		Rdb = redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "", // no password set
			DB:           0,  // use default DB
			WriteTimeout: 1000 * time.Second,
		})

	})
	return Rdb
}

func LoadFacToRedisHash(tblName string) {
	// PipeFacToRedis loads a fac file to redis using pipeline mode
	tblPath := FindFile(tblName)
	if !strings.HasSuffix(tblPath, ".fac") {
		// error
		panic("Not a fac file")

	} else {
		f, err := os.Open(tblPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		textlines := make([]string, 0)
		for scanner.Scan() {
			textlines = append(textlines, scanner.Text())
		}

		keyNames := make([]string, 0)
		hashKeys := make([]string, 0)
		noIdx := 0

		pipe := Rdb.Pipeline()
		for _, line := range textlines {

			line = strings.ReplaceAll(line, "\"", "")
			// Header
			if strings.HasPrefix(line, "!") {
				noIdx, err = strconv.Atoi(line[1:2])
				if err != nil {
					panic(err)
				}
				noKeys := noIdx - 1
				if noKeys < 1 {
					panic(tblName + " has no keys")
				}

				// remove ""
				str := strings.Split(line, ",")
				keyNames = str[1:noIdx]
				hashKeys = str[noIdx:]

				pipe.RPush(ctx, tblName+":000Key1", keyNames)
				pipe.RPush(ctx, tblName+":000Key2", hashKeys)
			}
			// Records
			if strings.HasPrefix(line, "*") {
				//process records
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				hashValues := str[noIdx:]
				for i, v := range hashValues {
					pipe.HSet(ctx, tblName+":"+key, hashKeys[i], v)
				}
			}
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}
	}
}
