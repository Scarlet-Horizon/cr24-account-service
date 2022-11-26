package util

import (
	"github.com/google/uuid"
	"log"
)

const TableName = "Account"

var AccountTypesLimit = map[string]int{
	"checking": 50,
	"saving":   10,
}

func Log(message string, err error) {
	log.Println(message)
	log.Println(err)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetPK(id string) string {
	return "USER#" + id
}

func GetSK(id string) string {
	return "ACCOUNT#" + id
}
