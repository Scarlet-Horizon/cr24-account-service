package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"main/db"
	"main/env"
	"main/model"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env variables: %v", err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Fatalf("failed to load SDK config, %v", err)
	}

	acc := db.AccountDB{
		Client: dynamodb.NewFromConfig(cfg),
	}

	bankAccount := model.Account{
		PK:       "USER#4545",
		SK:       "ACCOUNT#" + uuid.NewString(),
		Amount:   100.52,
		Limit:    50,
		OpenDate: time.Now(),
		Type:     "main",
	}
	err = acc.Create(bankAccount)

	//err = acc.Create(util.RandomAccount())
	fmt.Println(err)
}
