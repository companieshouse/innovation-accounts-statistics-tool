package config

import (
	"errors"
	"fmt"
	"github.com/companieshouse/gofigure"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

// Config holds configuration details set by the environment.
type Config struct {
	TransactionsMongoDBURL      string `env:"TRANSACTIONS_MONGODB_URL"       flag:"transactions-mongodb-url"       flagDesc:"Transactions MongoDB server URL"`
	TransactionsMongoDBDatabase string `env:"TRANSACTIONS_MONGODB_DATABASE"  flag:"transactions-mongodb-database"  flagDesc:"Transactions MongoDB database for data"`
	LogLevel                    string `env:"LOG_LEVEL"                      flag:"log-level"                      flagDesc:"Logging level of the application"`
	SenderEmail                 string `env:"SENDER_EMAIL"                   flag:"sender-email"                   flagDesc:"Email of Sender"`
	ReceiverEmails              string `env:"RECEIVER_EMAILS"                flag:"receiver-emails"                flagDesc:"Emails of each Receiver"`
	SesAwsRegion                string `env:"SES_AWS_REGION"                 flag:"ses-aws-region"                 flagDesc:"AWS Region"`
}

var cfg *Config
var mtx sync.Mutex

// Get returns a pointer to a Config instance,
// populated with values from environment or command-line flags.
func Get() (*Config, error) {

	mtx.Lock()
	defer mtx.Unlock()

	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}

	mandatoryConfigsMissing := false

	if cfg.TransactionsMongoDBURL == "" {
		log.Warn("TRANSACTIONS_MONGODB_URL not set in environment")
		mandatoryConfigsMissing = true
	}

	if cfg.TransactionsMongoDBDatabase == "" {
		log.Warn("TRANSACTIONS_MONGODB_DATABASE not set in environment")
		mandatoryConfigsMissing = true
	}

	if cfg.SenderEmail == "" {
		log.Warn("SENDER_EMAIL not set in environment")
		mandatoryConfigsMissing = true
	}

	if cfg.ReceiverEmails == "" {
		log.Warn("RECEIVER_EMAILS not set in environment")
		mandatoryConfigsMissing = true
	}

	if cfg.SesAwsRegion == "" {
		log.Warn("SES_AWS_REGION not set in environment")
		mandatoryConfigsMissing = true
	}

	if mandatoryConfigsMissing {
		return nil, errors.New("mandatory configs missing from environment")
	}

	return cfg, nil
}

// SetLogLevel sets the level of logging using the given ENV variable "LOG_LEVEL".
func SetLogLevel(cfg *Config) {

	if cfg.LogLevel != "" {
		log.Info(fmt.Sprintf("Log level set in environment, attempting to set log level to: %s", cfg.LogLevel))
		lvl, err := log.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Error(fmt.Sprintf("failed to set log level: %s. Exiting", err))
			os.Exit(1)
		}
		log.SetLevel(lvl)
		log.Info("Log level set successfully")
	}
}
