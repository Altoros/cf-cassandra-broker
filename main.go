package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/Altoros/cf-cassandra-service-broker/app"
	"github.com/Altoros/cf-cassandra-service-broker/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")

	flag.Parse()
}

func main() {
	if configFile == "" {
		log.Fatal("No config file specified")
	}
	config, err := config.InitFromFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := app.NewApp(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	k := make(chan os.Signal, 1)
	signal.Notify(k, os.Interrupt, os.Kill)
	go func() {
		<-k
		app.Stop()
	}()
	app.Start()
}
