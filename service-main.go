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
func setConfig() (*daemon.Config, error) {
	cfg := daemon.Config{}

	var pathToConfigFile string

	flag.StringVar(&pathToConfigFile, "cfg", "./config.json", "Path to file with configuration for service in JSON format")
	flag.Parse()

	data, err := ioutil.ReadFile(pathToConfigFile)
	if err != nil {
		log.Println("Can't read config file", err.Error())
		return nil, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Println("Can't unmarshal config file", err.Error())
		return nil, err
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

	return &cfg, nil
}

func main() {
	log.SetFlags(log.Llongfile)
	cfg, err := setConfig()
	if err != nil {
		log.Fatalln("Error config", err.Error())
	}
	if err := daemon.RunService(cfg); err != nil {
		log.Fatalln("Error RunService", err.Error())
	}
}
