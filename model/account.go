package model

import (
	"encoding/json"
	"strings"
	"time"
)

func getUserID(id string) string {
	if strings.Contains(id, "USER#") {
		return strings.Split(id, "USER#")[1]
	}
	return id
}

func getAccountID(id string) string {
	if strings.Contains(id, "ACCOUNT#") {
		return strings.Split(id, "ACCOUNT#")[1]
	}
	return id
}

type Account struct {
	// User UUID
	PK string `dynamodbav:"PK" json:"userID" example:"6204037c-30e6-408b-8aaa-dd8219860b4b"`
	// Account UUID
	SK string `dynamodbav:"SK" json:"accountID" example:"09130407-1f81-4ac5-be85-6557683462d0"`
	// Account amount
	Amount float64 `dynamodbav:"Amount" json:"amount" example:"50.5"`
	// Account limit
	Limit int `dynamodbav:"Limit" json:"limit" example:"50"`
	// The opening date for the account
	OpenDate time.Time `dynamodbav:"OpenDate,unixtime" json:"openDate" example:"2022-11-26T11:59:38+01:00"`
	// The closing date for the account
	CloseDate *time.Time `dynamodbav:"CloseDate,omitempty,unixtime" json:"closeDate,omitempty" example:"2022-12-21T14:40:20+01:00"`
	// Account type. One of the following: 'checking', 'saving'
	Type string `dynamodbav:"Type" json:"type" example:"checking" enums:"checking,saving"`
	// Account transactions
	Transactions []Transaction `dynamodbav:"omitempty" json:"transactions,omitempty"`
} //@name Account

func (account Account) MarshalJSON() ([]byte, error) {
	type Alias Account
	return json.Marshal(&struct {
		PK string `dynamodbav:"PK" json:"userID" example:"6204037c-30e6-408b-8aaa-dd8219860b4b"`
		SK string `dynamodbav:"SK" json:"accountID" example:"09130407-1f81-4ac5-be85-6557683462d0"`
		*Alias
	}{
		PK:    getUserID(account.PK),
		SK:    getAccountID(account.SK),
		Alias: (*Alias)(&account),
	})
}
