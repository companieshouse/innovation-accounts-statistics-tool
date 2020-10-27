package service

import (
	"fmt"
	"github.com/companieshouse/innovation-accounts-statistics-tool/config"
	"github.com/companieshouse/innovation-accounts-statistics-tool/models"
	"os"
	"time"

	"github.com/companieshouse/innovation-accounts-statistics-tool/db"
	log "github.com/sirupsen/logrus"
)

const statisticsReportFileNamePrefix = "CHS_SmallFullAccounts_Statistics"
const csvFileSuffix = ".csv"

// Service provides an interface to retrieve statistics in a CSV format.
type Service interface {
	GetStatisticsReport(dataDescription string) *models.CSV
}

// Impl is a concrete implementation of the Service interface.
type Impl struct {
	transactionClient db.TransactionClient
}

// NewService returns a new service interface implementation.
func NewService(cfg *config.Config) Service {
	return &Impl{
		transactionClient: db.NewTransactionDatabaseClient(cfg),
	}
}

// GetStatisticsReport returns a CSV version of the StatisticsReport struct, which is comprised of transactions data.
func (s *Impl) GetStatisticsReport(dataDescription string) *models.CSV {

	transactions, err := s.transactionClient.GetAccountsTransactions(dataDescription)
	if err != nil {
		log.Error(fmt.Sprintf("Error when retrieving transactions: %s", err))
		os.Exit(1)
	}

	sr := sortTransactionsPerMonth(transactions)

	// Print will only log at trace level.
	printStatisticsReport(sr)

	csv := constructCSV(sr)

	return &csv
}

// sortTransactionsPerMonth takes a slice of Transaction and groups them by the following criteria:
// 1. Grouped by Status being either Accepted or Rejected.
// 2. Grouped by of filing.
// 3. Grouped by month of filing.
// Returns a StatisticsReport model.
func sortTransactionsPerMonth(transactions *[]models.Transaction) *models.StatisticsReport {

	sr := models.NewStatisticsReport()

	sr.ClosedTransactions = len(*transactions)

	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for _, t := range *transactions {
		accepted := t.Data.Filings[t.ID+"-1"].Status == "accepted"
		rejected := t.Data.Filings[t.ID+"-1"].Status == "rejected"

		if accepted {
			if t.Data.ClosedAt.After(oneYearAgo) {
				sr.FirstYearAcceptedMonthlyFilings[t.Data.ClosedAt.Month()]++
			}
			sr.AcceptedTransactions++

		} else if rejected {
			sr.RejectedTransactions++
		}
	}

	return sr
}

// printStatisticsReport logs at Trace level and prints the stats report.
func printStatisticsReport(sr *models.StatisticsReport) {

	// Filings for the first year, printed per month.
	log.Traceln(fmt.Sprintf("--- Statistics Report Tool ---"))
	log.Traceln(fmt.Sprintf("--- Within 12 months Filings (Per Month) ---"))

	for month, total := range sr.FirstYearAcceptedMonthlyFilings {
		log.Traceln(fmt.Sprintf("%v Filings: %d", month.String(), total))
	}

	log.Traceln(fmt.Sprintf("--- Total: %d ---", sr.ClosedTransactions))
	log.Traceln(fmt.Sprintf("-------------------"))

	// Total filings printed per status.
	log.Traceln(fmt.Sprintf("--- Filings grouped by status ---"))
	log.Traceln(fmt.Sprintf("Closed transactions: %d", sr.ClosedTransactions))
	log.Traceln(fmt.Sprintf("Accepted transactions: %d", sr.AcceptedTransactions))
	log.Traceln(fmt.Sprintf("Rejected transactions: %d", sr.RejectedTransactions))
	log.Traceln(fmt.Sprintf("-------------------"))
}

// constructCSV marshals CSVable data into a CSV, accompanied by a file name.
func constructCSV(data models.CSVable) models.CSV {

	return models.CSV{
		Data:     data,
		FileName: statisticsReportFileNamePrefix + csvFileSuffix,
	}
}
