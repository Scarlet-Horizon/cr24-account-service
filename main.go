package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"main/db"
	"main/env"
	"main/util"
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

	err = acc.Create(util.RandomAccount())
	fmt.Println(err)
}
