package models

import (
	"time"
)

// Transaction describes a transaction database entity.
type Transaction struct {
	ID   string          `bson:"_id"`
	Data TransactionData `bson:"data"`
}

// TransactionData holds the data of each Transaction.
type TransactionData struct {
	Description string                       `bson:"description"`
	Filings     map[string]TransactionFiling `bson:"filings"`
	Links       map[string]string            `bson:"links"`
	Status      string                       `bson:"status"`
	ClosedAt    time.Time                    `bson:"created_at"`
}

// TransactionFiling contains the type and status of a Transactions Filing.
type TransactionFiling struct {
	Type   string `bson:"type"`
	Status string `bson:"status"`
}
