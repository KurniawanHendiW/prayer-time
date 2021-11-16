package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type service struct {
	redisPool *redis.Pool
}

type Service interface {
	Get(key string) *Reply
	Set(key string, value interface{}) *Reply
	Del(key string) *Reply
	Expire(key string, expireIn int64) *Reply
}

func NewService(config RedisConfig) Service {
	s := &service{
		redisPool: initRedis(config),
	}

	if err := s.testRedis(); err != nil {
		log.Fatal(err)
	}

	return s
}

func initRedis(config RedisConfig) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(config.Timeout) * time.Second,
		Wait:        true,
		Dial: func() (c redis.Conn, err error) {
			if config.TlsUrl != "" {
				c, err = redis.DialURL(config.TlsUrl, redis.DialTLSSkipVerify(true))
				if err != nil {
					return nil, err
				}
			} else {
				c, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))
				if err != nil {
					return nil, err
				}

				if config.Password != "" {
					if _, err = c.Do("AUTH", config.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redisPool
}

func (s *service) testRedis() error {
	if err := s.Set("TEST", "TEST").Error; err != nil {
		return err
	}

	if err := s.Del("TEST").Error; err != nil {
		return err
	}

	return nil
}

func (s *service) getConnection() redis.Conn {
	if s.redisPool != nil {
		return s.redisPool.Get()
	}

	return nil
}

func (s *service) do(command string, args ...interface{}) *Reply {
	conn := s.getConnection()
	defer conn.Close()

	result, err := conn.Do(command, args...)

	return &Reply{
		Result: result,
		Error:  err,
	}
}

func (s *service) Get(key string) *Reply {
	return s.do("GET", key)
}

func (s *service) Set(key string, value interface{}) *Reply {
	return s.do("SET", key, value)
}

func (s *service) Del(key string) *Reply {
	return s.do("DEL", key)
}

func (s *service) Expire(key string, expireIn int64) *Reply {
	return s.do("EXPIRE", key, fmt.Sprintf("%d", expireIn))
}

func (rp *Reply) Unmarshal(object interface{}) error {
	str, err := redis.String(rp.Result, rp.Error)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(str), object)
	if err != nil {
		return err
	}

	return nil
}

func (rp *Reply) String() (string, error) {
	return redis.String(rp.Result, rp.Error)
}

func (rp *Reply) Int64() (int64, error) {
	return redis.Int64(rp.Result, rp.Error)
}

func (rp *Reply) ArrByte() ([]byte, error) {
	return redis.Bytes(rp.Result, rp.Error)
}

func IsErrorNil(err error) bool {
	return err == redis.ErrNil
}
