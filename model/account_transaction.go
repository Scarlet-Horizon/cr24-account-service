package model

type AccountTransaction struct {
	Account
	// An array of Transaction
	Transactions []Transaction `json:"transactions"`
} //@name AccountTransaction
