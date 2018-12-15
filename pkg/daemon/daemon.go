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

	"bmstu.codes/developers34/SBWeb/pkg/s3"

	"bmstu.codes/developers34/SBWeb/pkg/model"

	"bmstu.codes/developers34/SBWeb/pkg/api"
	"bmstu.codes/developers34/SBWeb/pkg/db"
	sm "bmstu.codes/developers34/SBWeb/pkg/sessionmanager"
)

// Config is config structure for whole service.
type Config struct {
	DB  db.Config
	SM  sm.Config
	API api.Config
	IM  s3.Config
}

// RunService is a function that starts the whole service using
// provided config. Every error that happens during this function
// is fatal. So program can't run if any error happens.
func RunService(cfg *Config) error {
	// init connection with database
	db, err := db.InitConnDB(cfg.DB)
	if err != nil {
		log.Println("Can't connect to database", err.Error())
		return err
	}
	log.Println("Connected to DB")

	// init connection with session manager
	sm, err := sm.InitConnSM(cfg.SM)
	if err != nil {
		log.Println("Can't start session manager", err.Error())
		return err
	}
	log.Println("Connected to SM")

	im, err := s3.InitS3(cfg.IM)
	if err != nil {
		log.Println("Can't access AWS S3", err.Error())
		return err
	}
	log.Println("Connected to AWS S3")

	// create model for API
	m := model.New(db, sm, im)

	// start server
	log.Println("Starting API server...")
	srv, ch := api.StartServer(cfg.API, m)

	// wait signal of server shutdown
	waitForSignal(srv, ch)

	return nil
}

// waitForSignal waits signal from OS to shutdown server and
// error from server himself.
func waitForSignal(srv *http.Server, chSrv chan error) error {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-ch:
			srv.Shutdown(nil)
			log.Printf("Got signal: %v, exiting.", s)
			return nil
		case err := <-chSrv:
			log.Println(err)
			return err
		}
	}
}
