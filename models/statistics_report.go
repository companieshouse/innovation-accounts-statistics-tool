package models

import (
	"strconv"
	t "time"
)

const (
	csvColumnRow   = 2  // 1 Header row + 1 Data row (2)
	csvColumnCount = 15 // 12 Months of the year + Closed, Accepted and Rejected (15)
)

// StatisticsReport holds statistical data formed from Transaction data.
type StatisticsReport struct {
	ClosedTransactions   int
	AcceptedTransactions int
	RejectedTransactions int

	FirstYearAcceptedMonthlyFilings  map[t.Month]int
	SecondYearAcceptedMonthlyFilings map[t.Month]int
}

// NewStatisticsReport returns a newly constructed StatisticsReport with default values.
func NewStatisticsReport() *StatisticsReport {
	return &StatisticsReport{
		ClosedTransactions:               0,
		AcceptedTransactions:             0,
		RejectedTransactions:             0,
		FirstYearAcceptedMonthlyFilings:  initialiseMap(),
		SecondYearAcceptedMonthlyFilings: initialiseMap(),
	}
}

// initialiseMap returns a map with months mapped to 0 values ready to be used.
func initialiseMap() map[t.Month]int {
	return map[t.Month]int{
		t.January:   0,
		t.February:  0,
		t.March:     0,
		t.April:     0,
		t.May:       0,
		t.June:      0,
		t.July:      0,
		t.August:    0,
		t.September: 0,
		t.October:   0,
		t.November:  0,
		t.December:  0,
	}
}

// ToCSV returns a [][]string version of the data within the StatisticsReport struct provided.
func (sr *StatisticsReport) ToCSV() [][]string {

	csv := make([][]string, csvColumnRow)

	csv[0] = sr.constructHeaders()
	csv[1] = sr.getValues()

	return csv
}

// constructHeaders retrieves the headers from the statistics report (Months of the year, and other important information)
// which will be used in the final CSV document as titles of each data point.
func (sr *StatisticsReport) constructHeaders() []string {

	headers := make([]string, csvColumnCount)

	for k := range sr.FirstYearAcceptedMonthlyFilings {
		headers[int(k)-1] = k.String()
	}

	headers[12] = "Total Closed"
	headers[13] = "Total Accepted"
	headers[14] = "Total Rejected"

	return headers
}

// getValues retrieves the data points which will sit under the previously retrieved headers in the new CSV file.
func (sr *StatisticsReport) getValues() []string {

	values := make([]string, csvColumnCount)

	counter := 0 // Use counter as Months start at 1, but we want our array to start at 0.
	for _, v := range sr.FirstYearAcceptedMonthlyFilings {
		values[counter] = strconv.Itoa(v)
		counter++
	}

	values[12] = strconv.Itoa(sr.ClosedTransactions)
	values[13] = strconv.Itoa(sr.AcceptedTransactions)
	values[14] = strconv.Itoa(sr.RejectedTransactions)

	return values
}
