package main

import (
	"flag"
	"os"

	"github.com/pressly/selfie"
	"github.com/pressly/selfie/config"
	"github.com/pressly/selfie/logme"
)

var (
	flags    = flag.NewFlagSet("selfie", flag.ExitOnError)
	confFile = flags.String("config", "", "path to config file")
)

func main() {
	flags.Parse(os.Args[1:])

	var err error
	var conf *config.Config

	//load configuration from either confFile or Env's CONFIG variable
	conf, err = config.New(*confFile, os.Getenv("CONFIG"))
	if err != nil {
		logme.Fatal(err)
	}

	//create a new Releasidier app.
	app, err := selfie.New(conf)
	if err != nil {
		logme.Fatal(err)
	}

	//start the selfie's App.
	//this will block until app stops, either by panic or exit signal
	app.Start()

	logme.Info("App is shutting down.")
	app.Exit()
}
