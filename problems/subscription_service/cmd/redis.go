package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/pandaemoniumplaza/concurrency/problems/subscription_service/data"
	"net/http"
	"os"
	"sync"
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

func initMailer(wg *sync.WaitGroup) *Mailer {
	return &Mailer{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromAddress: "info@mycompany.com",
		FromName:    "info",
		Wait:        wg,
		MailerChan:  make(chan Message, 100),
		ErrorChan:   make(chan error),
		DoneChan:    make(chan bool),
	}
}
