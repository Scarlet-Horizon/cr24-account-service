package model

import "time"

type Account struct {
	PK            string    `dynamodbav:"PK"`
	SK            string    `dynamodbav:"SK"`
	AccountNumber int       `dynamodbav:"AccountNumber"`
	Amount        float64   `dynamodbav:"Amount"`
	Limit         float64   `dynamodbav:"Limit"`
	OpenDate      time.Time `dynamodbav:"OpenDate"`
	CloseDate     time.Time `dynamodbav:"CloseDate"`
	Type          string    `dynamodbav:"Type"`
}
