package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var RDS *redis.Client

func InitRDS() {
	host := ""
	port := 6379
	pwd := ""
	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	RDS = rdb
}

func RPUSH() {

}
