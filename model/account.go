package model

import "time"

type Account struct {
	PK        string     `dynamodbav:"PK" json:"userID"`
	SK        string     `dynamodbav:"SK" json:"accountID"`
	Amount    float64    `dynamodbav:"Amount" json:"amount"`
	Limit     int        `dynamodbav:"Limit" json:"limit"`
	OpenDate  time.Time  `dynamodbav:"OpenDate,unixtime" json:"openDate"`
	CloseDate *time.Time `dynamodbav:"CloseDate,omitempty,unixtime" json:"closeDate,omitempty"`
	Type      string     `dynamodbav:"Type" json:"type"`
}
