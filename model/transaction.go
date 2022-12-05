package model

import (
	"time"
)

type Transaction struct {
	// Transaction UUID
	ID string `json:"id" example:"4a5ed2e0-5cdb-4f9e-96e3-ecc372ba4f0c"`
	// Sender account UUID
	SenderID string `json:"senderID" example:"5d84ca00-c079-4577-9560-e1014086affe"`
	// Recipient account UUID
	RecipientID string `json:"recipientID" example:"8cca0453-8e84-4f3b-aa40-7fc9cd162a34"`
	// Transaction amount
	Amount float64 `json:"amount" example:"17.24"`
	// Transaction date
	Date time.Time `json:"date" example:"2022-12-21T08:45:12+01:00"`
	// Transaction type
	Type TransactionType `json:"type"`
} //@name Transaction

func (receiver Transaction) GetDate() string {
	return receiver.Date.Format("2006-01-02 15:04:05")
}

type TransactionType struct {
	// TransactionType ID
	ID int `json:"id" example:"1"`
	// TransactionType description
	Type string `json:"type" example:"card-payment"`
} //@name TransactionType
