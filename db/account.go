package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"main/model"
	"main/util"
	"time"
)

type AccountDB struct {
	Client *dynamodb.Client
}

func (receiver AccountDB) Create(account model.Account) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"PK":       &types.AttributeValueMemberS{Value: account.PK},
			"SK":       &types.AttributeValueMemberS{Value: account.SK},
			"Amount":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", account.Amount)},
			"Limit":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", account.Limit)},
			"OpenDate": &types.AttributeValueMemberS{Value: account.OpenDate.Format("2006-01-02")},
			"Type":     &types.AttributeValueMemberS{Value: account.Type},
		},
		TableName: aws.String(util.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := receiver.Client.PutItem(ctx, input)
	return err
}
