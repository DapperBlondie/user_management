package main

import (
	"github.com/DapperBlondie/user_management/src/models"
	"github.com/DapperBlondie/user_management/src/repo"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const (
	PORT = ":8080"
	HOST = "localhost"
	PostgresDBString = "host=localhost port=5720 dbname=postgres user=postgres password=alireza1380##"
)

var appConfig *models.AppConfig
var scsManager *scs.SessionManager

func main()  {
	scsManager = scs.New()
	scsManager.Cookie.SameSite = http.SameSiteLaxMode
	scsManager.Cookie.Persist = true
	scsManager.Cookie.Secure = false
	scsManager.Lifetime = time.Hour*1

	appConfig = &models.AppConfig{SCSManager: scsManager}
	pgDB, err := repo.ConnectSQL(PostgresDBString)
	if err != nil {
		log.Fatalln("Error in Connecting To Database")
		return
	}

	NewHandlerRepo(appConfig, pgDB)

	srv := &http.Server{
		Addr:             	HOST+PORT,
		Handler:           routes(),
		ReadTimeout:       time.Second*8,
		ReadHeaderTimeout: time.Second*4,
		WriteTimeout:      time.Second*8,
		IdleTimeout:       time.Second*6,
	}

	log.Println("Listening and Serving on localhost:8080 ...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("Error in start the server")
		return
	}
	return
}
