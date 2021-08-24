package infrastructure

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type IRedisClient interface {
	Redis() *redis.Client
}

type RedisClient struct {
	log      *log.Logger
	host     []string
	password string
	db       int
}

func NewRedisClient(host []string, password string, db int, log *log.Logger) IRedisClient {
	return &RedisClient{host: host, password: password, db: db, log: log}
}

func (client *RedisClient) Redis() *redis.Client {
	ctx:=context.Background()
	cl := redis.NewClient(&redis.Options{
		Addr:         client.host[0],
		Password:     client.password,
		DB:           client.db,
		DialTimeout:  time.Duration(30) * time.Second,
		WriteTimeout: time.Duration(30) * time.Second,
		ReadTimeout:  time.Duration(30) * time.Second,
	})

	if _, err := cl.Ping(ctx).Result(); err != nil {
		panic(err)
	}
	client.log.Println("redis client connected successfully !")
	// cl.Publish()
	return cl
}
