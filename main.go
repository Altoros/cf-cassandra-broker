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

	go func() {
		app.Start()
	}()
	handleSignals()
	app.Stop()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
