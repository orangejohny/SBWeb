// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
Package daemon contains functions to start the whole service.
It uses own config type to configure connection with database,
session manager and run API server.
*/
package daemon

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/orangejohny/SBWeb/pkg/model"

	"github.com/orangejohny/SBWeb/pkg/api"
	"github.com/orangejohny/SBWeb/pkg/db"
	sm "github.com/orangejohny/SBWeb/pkg/sessionmanager"
)

// Config is config structure for whole service.
type Config struct {
	DB  db.Config
	SM  sm.Config
	API api.Config
}

// RunService is a function that starts the whole service using
// provided config. Every error that happens during this function
// is fatal. So program can't run if any error happens.
func RunService(cfg *Config) error {
	// if deployed to Heroku
	if os.Getenv("PORT") != "" {
		cfg.API.Address = ":" + os.Getenv("PORT")
	}
	if os.Getenv("DATABASE_URL") != "" {
		cfg.DB.DBAddress = os.Getenv("DATABASE_URL")
	}
	if os.Getenv("REDIS_URL") != "" {
		cfg.SM.DBAddress = os.Getenv("REDIS_URL")
	}

	// init connection with database
	db, err := db.InitConnDB(cfg.DB)
	if err != nil {
		log.Fatalln("Can't connect to database", err.Error())
		return err
	}
	log.Println("Connected to DB")

	// init connection with session manager
	sm, err := sm.InitConnSM(cfg.SM)
	if err != nil {
		log.Fatalln("Can't start session manager", err.Error())
		return err
	}
	log.Println("Connected to SM")

	// create model for API
	m := model.New(db, sm)

	// start server
	log.Println("Starting API server...")
	srv, ch := api.StartServer(cfg.API, m)

	// wait signal of server shutdown
	waitForSignal(srv, ch)

	return nil
}

// waitForSignal waits signal from OS to shutdown server and
// error from server himself.
func waitForSignal(srv *http.Server, chSrv chan error) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

LOOP:
	for {
		select {
		case s := <-ch:
			srv.Shutdown(nil)
			log.Printf("Got signal: %v, exiting.", s)
			<-chSrv
			break LOOP
		case err := <-chSrv:
			log.Fatalln(err)
		}
	}
}
