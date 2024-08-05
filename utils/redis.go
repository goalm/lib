package utils

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func init() {
	Rdb = GetRedisClient()
}

var Ctx = context.Background()
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
	Rdb.FlushDB(Ctx)
	return Rdb
}

func LoadFacToRedisHash(filePath string) {
	// PipeFacToRedis loads a fac file to redis using pipeline mode
	tblName, err := FilePathToName(filePath)
	if err != nil {
		panic(err)
	}
	if IsFac(filePath) {
		file, err := os.Open(filePath)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		pipe := Rdb.Pipeline()
		noIdx := 0
		hasKeys := make([]string, 0)
		for scanner.Scan() {
			line := scanner.Text()

			if line == "\xA0" || line == "" {
				break
			}
			// dump descriptions
			if line[0] != '!' && line[0] != '*' {
				continue
			}
			line = strings.ReplaceAll(line, "\"", "")
			// Header for Hash
			if line[0] == '!' {
				noIdx, err = strconv.Atoi(line[1:2])
				if err != nil {
					panic(err)
				}
				noRowKeys := noIdx - 1
				if noRowKeys < 1 {
					log.Printf("Table %s has no keys", tblName)
					return
				}
				fields := strings.Split(line, ",")
				hasKeys = fields[noIdx:]
			}
			// Records
			if line[0] == '*' {
				record := strings.Split(line, ",")
				rowKey := record[1]
				for _, v := range record[2:noIdx] {
					rowKey = rowKey + ":" + v
				}

				for i, v := range record[noIdx:] {
					pipe.HSet(Ctx, tblName+":"+rowKey, hasKeys[i], v)
				}
			}
		}
		_, err = pipe.Exec(Ctx)
		if err != nil {
			panic(err)
		}
	}
}
