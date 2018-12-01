package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"bmstu.codes/developers34/SBWeb/pkg/daemon"
)

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

	return &cfg
}

func main() {
	log.SetFlags(log.Llongfile)
	cfg := setConfig()
	if err := daemon.RunService(cfg); err != nil {
		log.Fatalln("Error RunService", err.Error())
	}
}
