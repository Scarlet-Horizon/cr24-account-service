package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"io"
	"log"
	"main/model"
	"main/response"
	"net/http"
	"os"
	"strings"
	"time"
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

func GetTransactions(accountID, token string) ([]model.Transaction, error) {
	req, err := http.NewRequest(http.MethodGet, "http://transaction-api:8085/api/v1/transaction/"+accountID+"/all", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return []model.Transaction{}, err
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
		return nil, errors.New("error: " + string(data))
	}

	var tr []model.Transaction
	if err := json.Unmarshal(data, &tr); err != nil {
		return nil, err
	}
	return tr, nil
}

func ValidateToken(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
		context.Abort()
		return
	}

	values := strings.Split(token, "Bearer ")
	if len(values) != 2 {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "token is not set properly"})
		context.Abort()
		return
	}

	token = values[1]

	to, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		context.Abort()
		return
	}

	if !to.Valid {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid token"})
		context.Abort()
		return
	}

	if claims, ok := to.Claims.(jwt.MapClaims); ok {
		if claims["sub"] == "" {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid id"})
			context.Abort()
			return
		}

		if claims["iat"] == "" || claims["exp"] == "" {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "iat or exp not set"})
			context.Abort()
			return
		}

		tokenIat := time.Unix(int64(claims["iat"].(float64)), 0)
		if tokenIat.After(time.Now()) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "iat can't be in the future"})
			context.Abort()
			return
		}

		tokenExp := time.Unix(int64(claims["exp"].(float64)), 0)
		if tokenExp.Before(time.Now()) {
			context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "expired token"})
			context.Abort()
			return
		}

		context.Set("ID", claims["sub"])
		context.Set("token", token)
		context.Next()
		return
	}
	context.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "invalid token"})
}
