package util

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"main/model"
	"net/http"
	"strings"
)

const TableName = "Account"

var AlreadyExists = errors.New("account with this type already exists")
var InsufficientFounds = errors.New("insufficient funds")
var InvalidAccount = errors.New("invalid account")
var OpenAccount = errors.New("account is not closed")

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

func GetTransactions(accountID string) ([]model.Transaction, error) {
	res, err := http.Get("http://transaction-api:8085/api/v1/transaction/" + accountID + "/all")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Printf("Close error: %s\n", err)
		}
	}(res.Body)

	if res.StatusCode == http.StatusNoContent {
		return []model.Transaction{}, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(data))
	}

	var tr []model.Transaction
	if err := json.Unmarshal(data, &tr); err != nil {
		return nil, errors.New("encode data error: " + err.Error())
	}
	return tr, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PATCH")
		c.Next()
	}
}
