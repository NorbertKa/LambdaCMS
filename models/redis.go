package db

import (
	"errors"
	"strconv"

	"github.com/NorbertKa/LambdaCMS/config"
	redis "gopkg.in/redis.v5"
)

const (
	ErrCantConnectToRedis string = "Cant connect to REDIS DB"
)

type Redis struct {
	*redis.Client
}

func NewRedis(conf *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + strconv.Itoa(conf.Redis.Port_Int()),
		Password: conf.Redis.Password, // no password set
		DB:       conf.Redis.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.New(ErrCantConnectToRedis)
	}
	redis := Redis{client}
	return &redis, nil
}
