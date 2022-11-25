package util

import "log"

const TableName = "Account"

var AccountTypesLimit = map[string]int{
	"checking": 50,
	"saving":   10,
}

func Log(message string, err error) {
	log.Println(message)
	log.Println(err)
}
