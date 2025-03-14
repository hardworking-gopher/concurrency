package main

import (
	"fmt"
	"github.com/pandaemoniumplaza/goroutines/subscription_service/data"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	envPgConnStr    = "PG_CONN_STR"
	envRedisAddress = "REDIS_ADDRESS"

	appPort          = "8080"
	redisNetworkTcp  = "tcp"
	pgConnRetries    = 10
	pgConnRetryDelay = 500 * time.Millisecond
)

func main() {
	wg := &sync.WaitGroup{}

	app := App{
		Wait:     wg,
		Session:  initSession(),
		Mailer:   initMailer(wg),
		Models:   data.New(initDB()),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.LstdFlags),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.LstdFlags),
	}

	app.InfoLog.Println("starting the server")

	go app.listenForMail()
	go app.listerForShutdown()

	app.serve()
}

func (a *App) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", appPort),
		Handler: a.routes(),
	}

	a.InfoLog.Printf("serving on port %s\n", appPort)

	if err := srv.ListenAndServe(); err != nil {
		a.ErrorLog.Fatal(err)
	}

	a.InfoLog.Println("server has shut down")
}

func (a *App) listerForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	a.InfoLog.Printf("received %s signal\n", <-quit)
	a.shutdown()
	a.InfoLog.Println("server has shut down")
	os.Exit(0)
}

func (a *App) shutdown() {
	a.InfoLog.Println("initiating graceful shutdown")

	a.InfoLog.Println("waiting for all the emails to be sent")

	a.Wait.Wait()

	a.InfoLog.Println("sending stop signal to email listener")
	a.Mailer.DoneChan <- true

	close(a.Mailer.MailerChan)
	close(a.Mailer.ErrorChan)
	close(a.Mailer.DoneChan)

	a.InfoLog.Println("closing channels")
}
