// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
SBWeb is the web-service for course work of developers34 team.
It's supposed to search and hire specialists of building industry.

Usage of application

Environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be specified.
If environment variable PORT is specified then its value will override value of config API address.
If environment variable REDIS_URL is specified then its value will override value of config SM DBAddress.
If environment variable DATABASE_URL is specified then its value will override value of config DB DBAddress.

To run application you need to specify the "cfg" parameter that receives
path to config file formatted as JSON.

Config has this structure:
  {
    "DB": {
      "DBAddress": <Address of postgres database (string)>,
      "MaxOpenConns": <Number of maximum open connections to database (int)>
    },
    "SM": {
      "DBAddress": <Address of redis storage (string)>,
      "TockenLength": <Length of tocken that will be used as session tocken (int)>,
      "ExpirationTime": <Expiration time of a session in seconds (int)>
    },
    "API": {
      "Address": <Port where the server will be started (string)>,
      "ReadTimeout": <Maximum duration for reading the entire request, including the body (string with postfix 's')>,
      "WriteTimeout": <Maximum duration before timing out writes of the response (string with postfix 's')>,
      "IdleTimeout": <Maximum amount of time to wait for the next request when keep-alives are enabled (string with postfix 's')>
    },
    "IM": {
      "Bucket": <Name of the AWS S3 bucket where to store images (string)>,
      "ACL": <Permissions to images uploaded by application (string)>,
      "Region": <Region of the AWS S3 bucket (string)>
    }
  }
*/
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
