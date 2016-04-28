package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/Altoros/cf-cassandra-broker/broker"
	"github.com/Altoros/cf-cassandra-broker/config"
)

var (
	configFile string
	pidFile    string
)

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.StringVar(&pidFile, "p", "", "Pid file")

	flag.Parse()
}

func main() {
	if configFile == "" {
		log.Fatal("No config file specified")
	}
	if pidFile != "" {
		writePid()
	}

	config, err := config.InitFromFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}

	broker, err := broker.New(config)
	if err != nil {
		log.Fatalf("Error creating broker: %s", err.Error())
	}

	go func() {
		broker.Start()
	}()
	handleSignals()
	broker.Stop()
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func writePid() {
	pid := strconv.Itoa(os.Getpid())
	f, err := os.Create(pidFile)
	if err != nil {
		log.Println("Warning: can not write pid:", err.Error())
	}
	defer f.Close()
	f.WriteString(pid)
}
