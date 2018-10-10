package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"./Configuration"
	"./Helper/Http"
	"./Middlewares"
	"./Services/API"
	"./Services/Monitoring"
)

var (
	gitHash   = ""
	buildDate = ""
)

func main() {
	// TODO en attendant de régler le problème des packages locaux non accessible avec -ldflags
	configuration.SetBuildInfo(gitHash, buildDate)

	// Création du logger de sortie
	logger := log.New(os.Stderr, "", log.LstdFlags)
	// Affichage des informations du logiciel
	logger.Print(configuration.String())

	// Parse application command line options
	configFile := flag.String("config", "", "Server configuration file path")
	flag.Parse()

	// Lecture de la configuration du service
	var err = errors.New("No configuration set")
	if *configFile != "" {
		err = configuration.ReadAndCreate(*configFile)
	} else {
		err = configuration.ReadSpringCloudConfig()
	}
	if err != nil {
		flag.Usage()
		logger.Fatal(err)
	}
	config := configuration.Get().Configuration.Data

	// Create services runner
	servicesRunner := helperhttp.Services{Logger: logger}

	// API endpoints
	apiRouter := apiService.NewRouter()
	apiRouter.Use(middlewares.NewPrometheus("api-klit").Handler)
	servicesRunner.Add("API", config.API, apiRouter)

	// Monitoring endpoints
	servicesRunner.Add("Monitoring", config.Monitoring, monitoringService.NewRouter())

	// Run all services
	servicesRunner.Run()
}
