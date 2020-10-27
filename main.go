package main

import (
	"fmt"
	lam "github.com/aws/aws-lambda-go/lambda"
	"github.com/companieshouse/innovation-accounts-statistics-tool/config"
	"github.com/companieshouse/innovation-accounts-statistics-tool/lambda"
	log "github.com/sirupsen/logrus"
	"os"
)

// Main entry point for application, will create config from ENV variables and start the lambda.
func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Error(fmt.Sprintf("Error when establishing config: %s", err))
		os.Exit(1)
	}

	config.SetLogLevel(cfg)

	l := lambda.New(cfg)
	lam.Start(l.Execute)
}
