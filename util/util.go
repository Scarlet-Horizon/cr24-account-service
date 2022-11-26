package util

import (
	"errors"
	"github.com/google/uuid"
	"log"
	"strings"
)

const TableName = "Account"

var AlreadyExists = errors.New("account with this type already exists")
var InsufficientFounds = errors.New("insufficient funds")

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
	if strings.Contains(id, "USER#") {
		return id
	}
	return "USER#" + id
}

func GetSK(id string) string {
	if strings.Contains(id, "ACCOUNT#") {
		return id
	}
	return "ACCOUNT#" + id
}
