package lambda

import (
	"github.com/companieshouse/innovation-accounts-statistics-tool/aws"
	"github.com/companieshouse/innovation-accounts-statistics-tool/config"
	"github.com/companieshouse/innovation-accounts-statistics-tool/service"
	log "github.com/sirupsen/logrus"
)

// Lambda facilitates the execution of Company Accounts statistics retrieval.
type Lambda struct {
	Service service.Service
}

type jsonBody struct{}

// New returns a new Lambda using the provided configs.
func New(cfg *config.Config) *Lambda {
	return &Lambda{
		Service: service.NewService(cfg),
	}
}

// Execute handles lambda execution.
func (lambda *Lambda) Execute(j *jsonBody) error {

	srCSV := lambda.Service.GetStatisticsReport("CIC report and full accounts")

	eg := aws.NewEmailGenerator()

	cfg, err := config.Get()
	if err != nil {
		log.Error(err.Error())
	}

	err = eg.GenerateEmail(srCSV, cfg)
	if err != nil {
		log.Error(err.Error())
	}

	return nil
}
