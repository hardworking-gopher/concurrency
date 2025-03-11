package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/pandaemoniumplaza/goroutines/subscription_service/data"
	"net/http"
	"os"
	"time"
)

func initSession() *scs.SessionManager {
	connStr := os.Getenv(envRedisAddress)
	gob.Register(data.User{})

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	session.Store = redisstore.New(initRedis(connStr))

	return session
}

func initRedis(connStr string) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 3,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(redisNetworkTcp, connStr)
		},
	}

	return redisPool
}
