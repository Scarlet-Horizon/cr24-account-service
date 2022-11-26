package model

import "time"

type Account struct {
	// User UUID. Has prefix USER#
	PK string `dynamodbav:"PK" json:"userID" example:"USER#6204037c-30e6-408b-8aaa-dd8219860b4b"`
	// Account UUID. Has prefix ACCOUNT#
	SK string `dynamodbav:"SK" json:"accountID" example:"ACCOUNT#09130407-1f81-4ac5-be85-6557683462d0"`
	// Account amount
	Amount float64 `dynamodbav:"Amount" json:"amount" example:"50.5"`
	// Account limit
	Limit int `dynamodbav:"Limit" json:"limit" example:"50"`
	// The opening date for the account
	OpenDate time.Time `dynamodbav:"OpenDate,unixtime" json:"openDate" example:"2022-11-26T11:59:38+01:00"`
	// The closing date for the account
	CloseDate *time.Time `dynamodbav:"CloseDate,omitempty,unixtime" json:"closeDate,omitempty" example:"2022-12-21T14:40:20+01:00"`
	// Account type. One of the following: 'checking', 'saving'
	Type string `dynamodbav:"Type" json:"type" example:"checking"`
}
