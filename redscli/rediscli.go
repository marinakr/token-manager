package redscli

import (
	//
	"github.com/go-redis/redis"
	"encoding/json"
	//
)

type RedisENV interface {
}

type DBCli struct {
	*redis.Client
}

func New(redisCreds interface{}) DBCli {
	data, err_json := json.Marshal(redisCreds)
	if err_json == nil {
		var result redis.Options
		json.Unmarshal(data, &result)
		client := redis.NewClient(&result)
		_, err_ping := client.Ping().Result()
		if err_ping != nil {
			panic("No connection to redis")
		} else {
			return client
		}
	} else {
		panic("Error in config file: redis")
	}
}

