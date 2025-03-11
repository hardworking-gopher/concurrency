package main

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/pandaemoniumplaza/goroutines/subscription_service/data"
	"log"
	"sync"
)

type App struct {
	Session  *scs.SessionManager
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
	Models   data.Models
}
