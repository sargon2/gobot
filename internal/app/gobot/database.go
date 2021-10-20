package gobot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Database struct {
	svc *dynamodb.DynamoDB
}

func NewDatabase() *Database {
	mySession := session.Must(session.NewSession())

	svc := dynamodb.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	return &Database{
		svc: svc,
	}
}

func (d *Database) Put(tablename string, data interface{}) bool {
	item, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		fmt.Printf("Error in Put: %v\n", err)
		return false
	}
	var putItemInput = &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tablename),
	}
	_, err = d.svc.PutItem(putItemInput)
	if err != nil {
		fmt.Printf("Error in Put: %v\n", err)
		return false
	}
	return true
}

func (d *Database) Get(tablename string, item interface{}, key string) bool {
	result, err := d.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		fmt.Printf("Error in Get: %v\n", err)
		return false
	}

	if result.Item == nil {
		return false
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Printf("Error in Get: %v\n", err)
		return false
	}

	return true
}

func (d *Database) Delete(tablename string, key string) bool {
	_, err := d.svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"Key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		fmt.Printf("Error in Delete: %v\n", err)
		return false
	}
	return true
}
