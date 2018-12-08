// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// SBWeb is the web-service for course work of developers34 team.
// It's supposed to search and hire specialists of building industry.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"bmstu.codes/developers34/SBWeb/pkg/daemon"
)

// setConfig parses the config provided in command options
func setConfig() *daemon.Config {
	cfg := daemon.Config{}

	var pathToConfigFile string

	flag.StringVar(&pathToConfigFile, "cfg", "./config.json", "Path to file with configuration for service in JSON format")
	flag.Parse()

	data, err := ioutil.ReadFile(pathToConfigFile)
	if err != nil {
		log.Fatalln("Can't read config file", err.Error())
		return nil
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalln("Can't unmarshal config file", err.Error())
		return nil
	}

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

	return &cfg
}

func main() {
	log.SetFlags(log.Llongfile)
	cfg := setConfig()
	if err := daemon.RunService(cfg); err != nil {
		log.Fatalln("Error RunService", err.Error())
	}
}
