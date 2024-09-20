package ddbproxy

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var svc *dynamodb.Client

func init() {
	// Load the shared AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create a DynamoDB client
	svc = dynamodb.NewFromConfig(cfg)
}

func QueryTable[T interface{}](tableName string, keyExpression string, hashKeyAttributeValues map[string]types.AttributeValue) []T {
	// Create the query input parameters
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		KeyConditionExpression:    aws.String(keyExpression),
		ExpressionAttributeValues: hashKeyAttributeValues,
	}

	// Query the table
	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to query table, %v", err)
	}

	var items []T
	err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		log.Fatalf("failed to unmarshal query result items, %v", err)
	}

	return items
}

func CreateItem(tableName string, record interface{}) error {
	av, err := attributevalue.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("failed to marshal GroceryItem: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v: %v", err, av)
	}

	return nil
}

func UpdateItem(tableName string, key map[string]types.AttributeValue, record interface{}, ignoreKeys []string) error {
	av, err := attributevalue.MarshalMap(record)
	if err != nil {
		return fmt.Errorf("failed to marshal GroceryItem: %v", err)
	}

	updateExpression := "SET"
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)
	first := true

	for k, v := range av {
		// This is a hack to not update hash/sort keys (this will throw in ddb)
		if slices.Contains(ignoreKeys, k) {
			continue
		}

		if !first {
			updateExpression += ","
		}
		first = false
		updateExpression += fmt.Sprintf(" #%s = :%s", k, k)
		expressionAttributeNames["#"+k] = k
		expressionAttributeValues[":"+k] = v
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err = svc.UpdateItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func DeleteItem(tableName string, key map[string]types.AttributeValue) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	_, err := svc.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func BatchDeleteItems(tableName string, keys []map[string]types.AttributeValue) error {
	writeReqs := make([]types.WriteRequest, len(keys))
	for index, key := range keys {
		writeReqs[index] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: key,
			},
		}
	}

	_, err := svc.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})

	return err
}
