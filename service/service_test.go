package service

import (
	"errors"
	m "github.com/companieshouse/innovation-accounts-statistics-tool/db/mocks"
	"github.com/companieshouse/innovation-accounts-statistics-tool/models"
	"github.com/golang/mock/gomock"
	c "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	description = "description"
)

func TestImpl_GetStatisticsReport(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mClient := m.NewMockTransactionClient(mockCtrl)

	// ---- Error path: client returns error when getting transactions --- \\
	testTransactionsFromClientReturnsError(t, mClient)

	// ---- Success path: transactions are sorted and statistics struct is returned --- \\
	testTransactionsSortsCorrectly(t, createTransactionsSlice(), mClient)
}

func testTransactionsFromClientReturnsError(t *testing.T, mClient *m.MockTransactionClient) {
	c.Convey("Given I want a list of closed transactions which are sorted into a StatisticsReport struct", t, func() {
		svc := &Impl{
			transactionClient: mClient,
		}

		c.Convey("When I call the mongoDb client to get the transactions", func() {
			errorReturnedFromClient := errors.New("error getting transactions")
			mClient.EXPECT().GetAccountsTransactions(description).Return(nil, errorReturnedFromClient)

			c.Convey("Then an error is returned", func() {
				transactions, err := svc.transactionClient.GetAccountsTransactions(description)

				c.So(transactions, c.ShouldEqual, nil)
				c.So(err, c.ShouldEqual, errorReturnedFromClient)
			})
		})
	})
}

func testTransactionsSortsCorrectly(t *testing.T, transactions *[]models.Transaction, mClient *m.MockTransactionClient) {
	c.Convey("Given I want a list of closed transactions which are sorted into a StatisticsReport struct", t, func() {
		svc := &Impl{
			transactionClient: mClient,
		}

		c.Convey("When I pass my transactions into the sorting function", func() {

			mClient.EXPECT().GetAccountsTransactions(description).Return(transactions, nil)

			c.Convey("Then I am returned a valid CSV struct containing StatisticsReport data", func() {

				csv := svc.GetStatisticsReport(description)
				expectedCSVFileName := "CHS_SmallFullAccounts_Statistics.csv"

				c.So(csv.Data, c.ShouldNotBeNil)
				c.So(csv.FileName, c.ShouldEqual, expectedCSVFileName)
			})
		})
	})

}

func createTransactionsSlice() *[]models.Transaction {

	// filings date which will always be February just passed. (within a year).
	filingDate := time.Date(time.Now().AddDate(0, -5, 0).Year(), 2, 1, 0, 0, 0, 0, time.UTC)

	// Initialise empty slice of transactions.
	transactions := make([]models.Transaction, 0)

	// Create Accepted example data.
	filingAcceptedData := &models.TransactionFiling{
		Type:   "type",
		Status: "accepted",
	}
	filingAcceptedMap := make(map[string]models.TransactionFiling, 0)
	filingAcceptedMap["id0-1"] = *filingAcceptedData

	// Create Rejected example data.
	filingRejectedData := &models.TransactionFiling{
		Type:   "type",
		Status: "rejected",
	}
	filingRejectedMap := make(map[string]models.TransactionFiling, 0)
	filingRejectedMap["id1-1"] = *filingRejectedData

	// Dummy Accepted transaction.
	transactionDataAccepted := models.TransactionData{
		Description: description,
		Filings:     filingAcceptedMap,
		Links:       nil,
		Status:      "closed",
		ClosedAt:    filingDate,
	}

	// Dummy Rejected transaction.
	transactionDataRejected := models.TransactionData{
		Description: description,
		Filings:     filingRejectedMap,
		Links:       nil,
		Status:      "closed",
		ClosedAt:    filingDate,
	}

	// Add dummy data too transactions slice.
	transactions = append(transactions, models.Transaction{
		ID:   "id0",
		Data: transactionDataAccepted,
	}, models.Transaction{
		ID:   "id1",
		Data: transactionDataRejected,
	})

	// Return the transactions slice.
	return &transactions
}
