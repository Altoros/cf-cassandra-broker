package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Altoros/cf-cassandra-broker/config"
	"github.com/Altoros/cf-cassandra-broker/migrate"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration File")

	flag.Parse()
}

func main() {
	if configFile == "" {
		fmt.Fprintln(os.Stderr, "Error: no config file specified")
		flag.Usage()
		os.Exit(1)
	}
	config, err := config.InitFromFile(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config: "+err.Error())
		os.Exit(1)
	}

	err = migrate.Run(&config.Cassandra)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error migrating cassandra: "+err.Error())
		os.Exit(1)
	}
}
