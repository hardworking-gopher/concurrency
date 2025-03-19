package main

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/pandaemoniumplaza/goroutines/subscription_service/data"
	"log"
	"sync"
)

type App struct {
	Session       *scs.SessionManager
	DB            *sql.DB
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Wait          *sync.WaitGroup
	Models        data.Models
	Mailer        *Mailer
	ErrorChan     chan error
	ErrorChanDone chan bool
}

func (a *App) listenForMail() {
	a.InfoLog.Println("starting listening for emails")

	for {
		select {
		case msg := <-a.Mailer.MailerChan:
			a.Wait.Add(1)
			go a.Mailer.sendMail(msg, a.Mailer.ErrorChan)
		case err := <-a.Mailer.ErrorChan:
			a.ErrorLog.Println(err)
		case <-a.Mailer.DoneChan:
			a.InfoLog.Println("finished listening for emails")
			return
		}
	}
}

func (a *App) sendEmail(msg Message) {
	select {
	case a.Mailer.MailerChan <- msg:
	default:
		a.ErrorLog.Println("can't send email due to filled channel")
	}
}
