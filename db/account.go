package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"main/model"
	"main/util"
	"time"
)

type AccountDB struct {
	Client *dynamodb.Client
}

func (receiver AccountDB) Create(account model.Account) error {
	accItem, err := attributevalue.MarshalMap(account)
	if err != nil {
		return err
	}

	accInput := &dynamodb.PutItemInput{
		Item:      accItem,
		TableName: aws.String(util.TableName),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = receiver.Client.PutItem(ctx, accInput)
	return err
}

func (receiver AccountDB) GetAll(id string) ([]model.Account, error) {
	keyCond := expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(util.GetPK(id))),
		expression.Key("SK").BeginsWith("ACCOUNT#"),
	)

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(util.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := receiver.Client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var accounts []model.Account
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}
