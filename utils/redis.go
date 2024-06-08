package utils

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func LoadFacToRedis(tblPath, tblName string, rdb *redis.Client) {
	// LoadFacToRedis loads a fac file to redis
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

		rowIdxes := make([]string, 0)
		colIdxes := make([]string, 0)
		noIdx := 0

		for _, line := range textlines {
			// line
			// Header
			if strings.HasPrefix(line, "!") {
				noIdx, err = strconv.Atoi(line[1:2])
				if err != nil {
					panic(err)
				}

				str := strings.Split(line, ",")
				rowIdxes = str[1:noIdx]
				colIdxes = str[noIdx:]
				rdb.Set(ctx, tblName+":Idxes", strings.Join(rowIdxes, ", "), 0)
			}

			if strings.HasPrefix(line, "*") {
				//process records
				str := strings.Split(line, ",")
				rowKeys := str[1:noIdx]

				key := rowKeys[0]
				for _, v := range rowKeys[1:] {
					key = key + ":" + v
				}

				colKeys := str[noIdx:]
				for i, v := range colKeys {
					rdb.HSet(ctx, tblName+":"+key, colIdxes[i], v)
				}
			}
		}

	}
}
